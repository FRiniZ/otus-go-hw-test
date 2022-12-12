package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

const Usage = "Usage: \n%s --timeout=10s host port\n"

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout")
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		fmt.Printf(Usage, os.Args[0])
		os.Exit(1)
	}

	tcpclient := NewTelnetClient(net.JoinHostPort(flag.Args()[0], flag.Args()[1]), timeout, os.Stdin, os.Stdout)
	if err := tcpclient.Connect(); err != nil {
		fmt.Println("Can't connect")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		for {
			if err := tcpclient.Send(); err != nil {
				fmt.Println(err)
				stop()
				return
			}
		}
	}()

	go func() {
		for {
			if err := tcpclient.Receive(); err != nil {
				fmt.Println(err)
				stop()
				return
			}
		}
	}()

	<-ctx.Done()
}
