package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
)

func main() {

	// Check the consistency of input parameters ( host & port must be present, timeout must be in correct format, port must be a number)
	if err := validateArgs(nil); err != nil {
		flag.Usage()
		log.Fatal(err)
	}

	ctx := context.Background()

	client := NewClient(ctx, net.JoinHostPort(Host, Port), "tcp", duration, os.Stdin, os.Stdout)

	client.TerminateHandler()

	if err := client.Connect(); err != nil {
		log.Println("Connection closed with ", err)
		os.Exit(2)
	}
}
