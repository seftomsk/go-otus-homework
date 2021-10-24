package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "connection timeout")
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatalln("you must at least provide host and port: go-telnet host port")
	}
	if len(args) > 3 {
		log.Fatalln("invalid params list: " +
			strings.Join(args, ", ") +
			". Valid call with params go-telnet host port OR go-telnet timeout=10s host port")
	}
	var address string
	if len(args) == 2 {
		address = net.JoinHostPort(args[0], args[1])
	}
	if len(args) == 3 {
		address = net.JoinHostPort(args[1], args[2])
	}

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	var err error
	if err = client.Connect(); err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err = client.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	log.Println("...Connected to " + address)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		defer cancel()
		if err = client.Send(); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		defer cancel()
		if err = client.Receive(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-ctx.Done()
}
