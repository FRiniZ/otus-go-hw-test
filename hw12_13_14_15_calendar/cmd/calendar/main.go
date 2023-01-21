//go:generate protoc --go_out=../../api --proto_path=../../api/ ../../api/EventService.proto
//go:generate protoc --go-grpc_out=../../api --proto_path=../../api/ ../../api/EventServiceInterface.proto

package main

import (
	"fmt"
	"os"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/app"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/logger"
	internalgrpc "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/server/grpcservice"
	internalhttp "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/server/http"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
)

func main() {
	conf := NewConfig()
	storage := storage.NewStorage(conf.Storage.DB, conf.Storage.DSN)
	logger := logger.NewLogger(conf.Logger.Level, os.Stdout)
	calendar := app.NewCalendar(logger, conf.CalendarConf, storage)
	httpsrv := internalhttp.NewServer(logger, calendar, conf.HTTP.Host, conf.HTTP.Port)
	grpcsrv, _ := internalgrpc.NewServer(logger, calendar, conf.GRPC.Host, conf.GRPC.Port)

	calendar.Run(httpsrv, grpcsrv)

	fmt.Println("calendar stopped")
}
