package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/app"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/logger"
	internalhttp "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/server/http"
	memorystorage "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage/memory"
	sqlstorage "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	var idb app.Storage

	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()
	if err := config.LoadFileTOML(configFile); err != nil {
		fmt.Fprintf(os.Stderr, "Can't load config file:%v error: %v\n", configFile, err)
		os.Exit(1)
	}

	fmt.Println("Config:", config)

	logg, err := logger.New(config.Logger.Level, os.Stdin, nil)
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
			log.Fatalf("Can't connect to storage:%v\n", err) //nolint
		}
		idb = db
	}

	calendar := app.New(logg, idb)

	server := internalhttp.NewServer(logg, calendar, config.HTTPServer)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Errorf("failed to stop http server:%v\n", err)
		}
		if err := idb.Close(ctx); err != nil {
			logg.Errorf("failed to close db:%v\n", err)
		}
	}()

	logg.Infof("calendar is running...\n")

	if err := server.Start(ctx); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logg.Errorf("failed to start http server: %v\n", err.Error())
			cancel()
			os.Exit(1)
		}
	}
}
