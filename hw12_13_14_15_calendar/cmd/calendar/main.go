package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FRiniZ/otus-go-hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/FRiniZ/otus-go-hw-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/FRiniZ/otus-go-hw-test/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/FRiniZ/otus-go-hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/FRiniZ/otus-go-hw-test/hw12_13_14_15_calendar/internal/storage/sql"
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
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Config:", config)

	logg := logger.New(config.Logger.Level, os.Stdin, nil)

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
			panic(err)
		}
		idb = db
	}

	calendar := app.New(logg, idb)

	server := internalhttp.NewServer(logg, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Errorf("failed to stop http server: " + err.Error())
		}
	}()

	logg.Infof("calendar is running...\n")
	addr := net.JoinHostPort(config.HTTPServer.Host, config.HTTPServer.Port)
	if err := server.Start(ctx, addr); err != nil {
		logg.Errorf("failed to start http server: %v\n", err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
