//go:generate protoc --go_out=../../api --proto_path=../../api/ ../../api/EventService.proto
//go:generate protoc --go-grpc_out=../../api --proto_path=../../api/ ../../api/EventServiceInterface.proto

package main

import (
	"fmt"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/app"
)

func main() {
	config := NewConfig()

	calendar := app.NewCalendar(config.CalendarConf)
	calendar.Run()

	fmt.Println("calendar stopped")
}
