package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	Host    string
	Port    string
	Timeout uint
)

func init() {
	flag.UintVar(&Timeout, "timeout", 10, "Timeout to connect the destination host")
}

func validateArgs() {
	flag.Parse()

	fmt.Println(flag.NArg())
	if flag.NArg() == 0 || flag.NArg() < 2 {
		flag.Usage()
		log.Println("The host name & port must be specified")
		os.Exit(1)
	}

	for _, v := range flag.Args() {
		fmt.Println(v)
	}
	Host = flag.Arg(0)
	Port = flag.Arg(1)

	fmt.Printf("Host=%s, Port=%s, Timeout=%v\n", Host, Port, Timeout)

}
