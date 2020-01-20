package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
)

func main() {

	// Check the consistency of input parameters ( host & port must be present)
	if err := validateArgs(os.Args); err != nil {
		flag.Usage()
		log.Fatal(err)
	}

	// Parse the optional parameter Timeout in order to get the appropriate duration

	ctx := context.Background()

	client := NewClient(ctx, net.JoinHostPort(Host, Port), "tcp", duration, os.Stdin, os.Stdout)

	client.TerminateHandler()

	if err := client.Connect(); err != nil {
		log.Println("Error occured during the connect ", err)
	}
}
