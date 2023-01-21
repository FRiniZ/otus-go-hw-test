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
	config := NewConfig()
	storage := storage.NewStorage(config.Storage)
	logger := logger.NewLogger(config.Logger.Level, os.Stdout)
	calendar := app.NewCalendar(logger, config.CalendarConf, storage)
	httpsrv := internalhttp.NewServer(logger, calendar, config.HTTP.Host, config.HTTP.Port)
	grpcsrv, _ := internalgrpc.NewServer(logger, calendar, config.GRPC.Host, config.GRPC.Port)

	calendar.Run(httpsrv, grpcsrv)

	fmt.Println("calendar stopped")
}
