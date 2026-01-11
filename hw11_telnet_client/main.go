package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "set timeout")

	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Println(`no host or port specified. 
  usage: go-telnet [--timeout] <host> <port>`)
		os.Exit(1)
	}

	host, port := flag.Args()[0], flag.Args()[1]

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		slog.Error("error connecting from client", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	go func() {
		if err := client.Send(); err != nil {
			slog.Error("error sending data", "error", err)
			os.Exit(1)
		}
	}()

	go func() {
		if err := client.Receive(); err != nil {
			slog.Error("error receiving data", "error", err)
			os.Exit(1)
		}
	}()

	// wait until interrupt signal
	<-sigCh
}
