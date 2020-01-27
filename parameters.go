package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	// Host is ip or address or name
	Host string
	// Port is a destination port
	Port string
	// Timeout desired timeout for connection to remote end
	Timeout string

	duration time.Duration

	// ErrHostnPortMustBeDefined is a error to manage the case when Host or Port are not defined
	ErrHostnPortMustBeDefined error = errors.New("Host & port must be specified")
	// ErrPortMustBeANumber is error to catch the case when received port value is not a number
	ErrPortMustBeANumber error = errors.New("Port must be a number")
	// ErrInvalidTimeFormat is a error to catch the case when  invalid time fomrat is used for timeout parameter
	ErrInvalidTimeFormat error = errors.New("Invalid Time format. Time unit must be specified i.e. s,m,h etc")
)

func init() {

	flag.StringVar(&Timeout, "timeout", "10s", "Timeout to connect the destination host")
}

func validateArgs(args []string) error {

	flag.Usage = func() {
		filename := filepath.Base(os.Args[0])
		log.Println("Usage:", filename, "[--timeout=<value>] host port")
	}

	a := os.Args[1:]
	if args != nil {
		a = args
	}

	err := flag.CommandLine.Parse(a)

	if err != nil {
		return ErrInvalidTimeFormat
	}

	if flag.NArg() == 0 || flag.NArg() < 2 {
		return ErrHostnPortMustBeDefined
	}

	Host = flag.Arg(0)

	if _, err := strconv.Atoi(flag.Arg(1)); err != nil {
		return ErrPortMustBeANumber
	}
	Port = flag.Arg(1)

	duration, err = time.ParseDuration(Timeout)
	if err != nil {
		return ErrInvalidTimeFormat
	}
	return nil
}
