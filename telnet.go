package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// TelnetClient is a struct to store the information related to telnet client instance
type TelnetClient struct {
	ctxTimeOut        context.Context
	ctxCancel         context.Context
	proto             string        // "tcp" or "udp" TCP is used by default.
	destination       string        // host:port string
	timeout           time.Duration // related to command line timeout parameter
	connection        net.Conn
	cancelFuncTimeOut context.CancelFunc // Cancel Function for Context with Timeout
	cancelFunc        context.CancelFunc // Cancel Function for Context with Cancel
	in                io.Reader          // Reading from remote end
	out               io.Writer          // Writing to remote end
}

var (
	// ErrAbortByContext is an error to catch the case the the cancel() shooted.
	ErrAbortByContext error = errors.New("Aborted by context")
	//ErrAbortBySystem  error = errors.New("Aborted by system interrupt")
)

// NewClient creates a new instance of telnet client
func NewClient(ctx context.Context, dest string, protocol string, timeout time.Duration, in io.Reader, out io.Writer) *TelnetClient {

	if protocol == "" {
		protocol = "tcp"
	}

	tc := &TelnetClient{
		proto:       protocol,
		destination: dest,
		timeout:     timeout,
		in:          in,
		out:         out,
	}

	tc.ctxTimeOut, tc.cancelFuncTimeOut = context.WithTimeout(ctx, tc.timeout)
	tc.ctxCancel, tc.cancelFunc = context.WithCancel(ctx)

	return tc
}

// Connect make a connection to the remote end
func (t *TelnetClient) Connect() error {

	var err, ErrOut, ErrIn error

	dialer := &net.Dialer{}

	t.connection, err = dialer.DialContext(t.ctxTimeOut, t.proto, t.destination)

	if err != nil {
		return err
	}
	defer t.CloseSession()
	defer t.cancelFuncTimeOut()
	defer t.cancelFunc()

	log.Printf("Connected to host: %s with timeout=%v", t.destination, t.timeout)
	log.Println("Press Ctrl-D / Ctrl-C to exit")

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		// Read from t.connection (remote end) and write to t.out(os.Stdout)
		ErrOut = t.readerWriter(t.connection, t.out)
		wg.Done()
	}()

	go func() {
		// Read from t.in (os.Stdin) and write to t.connection(remote end)
		ErrIn = t.readerWriter(t.in, t.connection)
		wg.Done()
	}()

	wg.Wait()

	if ErrIn != nil && ErrOut != nil {
		if ErrIn == ErrOut {
			return fmt.Errorf(">%v", ErrIn)
		}
		return fmt.Errorf("In: %v, Out: %v", ErrIn, ErrOut)
	} else if ErrIn != nil {
		return ErrIn
	} else if ErrOut != nil {
		return ErrOut
	}

	return nil
}

func rw(ctx context.Context, cancel context.CancelFunc, in io.Reader, out io.Writer) error {
	var err error
	rw := bufio.NewScanner(in)
loop:
	for {
		select {
		case <-ctx.Done():
			err = ErrAbortByContext
			break loop
		default:
			if ok := rw.Scan(); ok {
				out.Write([]byte(rw.Text() + "\n"))
			} else {
				log.Println("<<<EOF>>> detected, aborting...Press Enter to close session")
				err = io.EOF
				cancel()
				break loop
			}
		}
	}
	return err
}

// Wrapper method for rw function
func (t *TelnetClient) readerWriter(in io.Reader, out io.Writer) error {

	return rw(t.ctxCancel, t.cancelFunc, in, out)

}

// CloseSession does a close of current session
func (t *TelnetClient) CloseSession() error {

	log.Println("Closing the session...")
	err := t.connection.Close()
	if nil != err {
		return err
	}
	return nil
}

// TerminateHandler catch an OS interrupts signals
func (t *TelnetClient) TerminateHandler() {
	go func() {
		terminate := make(chan os.Signal, 1)
		signal.Notify(terminate, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
		<-terminate
		log.Printf("Received Signal:%v aborting... Press Enter to close.", terminate)
		t.cancelFunc()
	}()
}
