package main

import (
	"fmt"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/app"
)

func main() {

	config := NewConfig()
	scheduler := app.NewScheduler(config.SchedulerConf)

	scheduler.Run()

	fmt.Println("calendar_scheduler stopped")
}
