package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// TelnetClient is a struct to store the information related to telnet client instance
type TelnetClient struct {
	destination string        // host:port string
	timeout     time.Duration // related to command line timeout parameter
	connection  net.Conn      // Connection to remote end
	in          io.Reader     // Reading from remote end
	out         io.Writer     // Writing to remote end
}

var (

	// ErrAbortedBySystemInterrupt is a error to indicate the program termination by Ctrl-C
	ErrAbortedBySystemInterrupt error = errors.New("Aborted by system interrupt")
)

// NewClient creates a new instance of TelnetClient
func NewClient(dest string, timeout time.Duration, in io.Reader, out io.Writer) *TelnetClient {

	tc := &TelnetClient{
		destination: dest,
		timeout:     timeout,
		in:          in,
		out:         out,
	}

	return tc
}

func (t *TelnetClient) connect() error {

	var err error
	
	t.connection, err = net.DialTimeout("tcp", t.destination, t.timeout)
	if err != nil {
		return err
	}
	// Don't forget to close the connection at the end.
	defer t.connection.Close()

	log.Printf("Connected to host: %s with timeout=%v", t.destination, t.timeout)
	log.Println("Press Ctrl-D / Ctrl-C to exit")

	ErrorChannel := make(chan error)
	InputChannel := make(chan string)
	OutputChannel := make(chan string)

	// Run the support goroutine to handle system interrupt with Ctrl-C
	go func() {
		sysInterrupt := make(chan os.Signal, 1)
		signal.Notify(sysInterrupt, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
		<-sysInterrupt
		log.Printf("Received system interrupt -%v, aborting...", sysInterrupt)
		ErrorChannel <- ErrAbortedBySystemInterrupt
	}()

	// Run reader for remote connection
	go readFromIO(t.connection, InputChannel, OutputChannel, ErrorChannel, true)

	// Run reader for stdin
	go readFromIO(t.in, InputChannel, OutputChannel, ErrorChannel, false)

	// Supervisor loop, stops in case of any error from the ErrorChannel.
	for {
		select {
		case in := <-InputChannel:
			fmt.Print(in)
		case out := <-OutputChannel:
			t.connection.Write([]byte(out))
		case err := <-ErrorChannel:
			return err
		}
	}

}

// readFromIO reads from the io.Reader and forward the receiced data in accordance to the direction flag
// direction = true  in(io.Reader) -> input (InputChannel)
// direction = false in(io.Reader) -> output (OutputChannel)
// In case of EOF or any other error - forward receiced error from io.Reader to the errchan (ErrorChannel)
func readFromIO(in io.Reader, input, output chan string, errchan chan error, direction bool) {
	r := bufio.NewReader(in)
	for {
		str, err := r.ReadString('\n')
		if err != nil {
			errchan <- err
		}
		if direction {
			input <- str
		} else {
			output <- str
		}
	}
}
