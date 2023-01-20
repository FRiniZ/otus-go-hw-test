//go:generate protoc --go_out=../../api --proto_path=../../api/ ../../api/EventService.proto
//go:generate protoc --go-grpc_out=../../api --proto_path=../../api/ ../../api/EventServiceInterface.proto

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/app"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/logger"
	grpcservice "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/server/grpcservice"
	internalhttp "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/server/http"
	memorystorage "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage/memory"
	sqlstorage "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage/sql"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/calendar_config.toml", "Path to configuration file")
}

func main() {
	var idb app.Storage

	flag.Parse()

	if flag.Arg(0) == "version" {
		app.PrintVersion()
		return
	}

	config := NewConfig()
	if err := config.LoadFileTOML(configFile); err != nil {
		fmt.Fprintf(os.Stderr, "Can't load config file:%v error: %v\n", configFile, err)
		os.Exit(1)
	}

	fmt.Println("Config:", config)

	log, err := logger.New(config.Logger.Level, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't allocate logger:%v\n", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	switch config.Storage.DB {
	case "in-memory":
		db := memorystorage.New()
		idb = db
	case "sql":
		db := sqlstorage.New(config.Storage.DSN)
		err := db.Connect(ctx)
		if err != nil {
			log.Fatalf("Can't connect to storage:%v\n", err) //nolint:gocritic
		}
		idb = db
	}

	calendar := app.New(log, idb)
	httpsrv := internalhttp.New(log, calendar, config.HTTPServer, cancel)

	unarayLoggerEnricherIntercepter := func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) { //nolint:gofumpt
		var b strings.Builder
		ip, _ := peer.FromContext(ctx)
		md, ok := metadata.FromIncomingContext(ctx)
		userAgent := "unknown"

		if ok {
			userAgent = md["user-agent"][0]
		}

		b.WriteString(ip.Addr.String())
		b.WriteString(" ")
		b.WriteString(time.Now().Format("02/Jan/2006:15:04:05 -0700"))
		b.WriteString(" ")
		b.WriteString(info.FullMethod)
		b.WriteString(" ")
		b.WriteString(userAgent)
		b.WriteString("\"\n")
		log.Infof(b.String())

		return handler(ctx, req)
	}

	basesrv := grpc.NewServer(grpc.UnaryInterceptor(unarayLoggerEnricherIntercepter))
	grpcsrv := grpcservice.New(log, calendar, config.GRPSServer, basesrv)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := grpcsrv.Stop(ctx); err != nil {
			if !errors.Is(err, grpc.ErrServerStopped) {
				log.Errorf("failed to stop GRPC-server:%v\n", err)
			}
		}

		if err := httpsrv.Stop(ctx); err != nil {
			if !errors.Is(err, http.ErrServerClosed) &&
				!errors.Is(err, context.Canceled) {
				log.Errorf("failed to stop HTTP-server:%v\n", err)
			}
		}
		if err := calendar.Close(ctx); err != nil {
			log.Errorf("failed to stop GRPC-server:%v\n", err)
		}

		if err := idb.Close(ctx); err != nil {
			log.Errorf("failed to close db:%v\n", err)
		}
	}()

	log.Infof("calendar is running...\n")

	g, ctxEG := errgroup.WithContext(ctx)
	func1 := func() error {
		return httpsrv.Start(ctxEG)
	}

	func2 := func() error {
		return grpcsrv.Start(ctxEG)
	}

	g.Go(func1)
	g.Go(func2)

	if err := g.Wait(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) &&
			!errors.Is(err, grpc.ErrServerStopped) &&
			!errors.Is(err, context.Canceled) {
			log.Errorf("%v\n", err)
		}
	}
}
