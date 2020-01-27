package main

import (
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

	client := NewClient(net.JoinHostPort(Host, Port), "tcp", duration, os.Stdin, os.Stdout)

	if err := client.connect(); err != nil {
		log.Println("Connection closed with ", err)
	}
}
