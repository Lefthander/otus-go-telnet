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
	cancelFuncTimeOut context.CancelFunc
	cancelFunc        context.CancelFunc
	in                io.Reader // Reading from remote end
	out               io.Writer // Writing to remote end
	err               error
	interruptChannel  chan bool
}

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

	var ErrOut, ErrIn error

	dialer := &net.Dialer{}

	t.connection, t.err = dialer.DialContext(t.ctxTimeOut, t.proto, t.destination)

	if t.err != nil {
		return t.err
	}
	defer t.connection.Close()
	defer t.cancelFuncTimeOut()
	defer t.cancelFunc()

	log.Println("Connected to host:", t.destination)
	log.Println("Press Ctrl-D / Ctrl-C to exit")

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		// Read fron t.in and write to t.conneciton
		ErrOut = t.readerWriter(t.connection, t.out)
		wg.Done()
		log.Println("Close conn -> t.out")
	}()
	//wg.Add(1)
	go func() {
		ErrIn = t.readerWriter(t.in, t.connection)
		wg.Done()
		log.Println("Close t.in -> conn")
	}()

	wg.Wait()

	t.CloseSession()

	if ErrIn != nil && ErrOut != nil {
		return fmt.Errorf("Error %v, %v", ErrIn, ErrOut)
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
			log.Println("Received ctx.Done()")
			err = errors.New("Aborted by context")
			break loop
		default:
			if ok := rw.Scan(); ok {
				out.Write([]byte(rw.Text() + "\n"))
			} else {
				log.Println("<<<EOF>>> detected, aborting...")
				err = io.EOF
				cancel()
				break loop
			}
		}
	}
	log.Println("Exit from the readerWriter()")
	return err
}

// Wrapper method for rw function
func (t *TelnetClient) readerWriter(in io.Reader, out io.Writer) error {

	t.err = rw(t.ctxCancel, t.cancelFunc, in, out)
	return t.err
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
		log.Println("Received Signal:", terminate)
		t.cancelFunc()
		t.interruptChannel <- true
	}()
}
