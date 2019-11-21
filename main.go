package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func readSocket(ctx context.Context, cancel context.CancelFunc, conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	for {
		select {
		case <-ctx.Done():
			break
		default:
			if !scanner.Scan() {
				log.Printf("Cannot read")
				cancel()
				break
			}
			text := scanner.Text()
			log.Printf("From server %s", text)
		}
	}
	log.Println("Finished readSocket")
}

func writeSocket(ctx context.Context, conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			break
		default:
			if !scanner.Scan() {
				break
			}
			s := scanner.Text()
			log.Println("To server", s)
			conn.Write([]byte(fmt.Sprintf("%s\n", s)))
		}
	}
	log.Println("Finished writeSocket")
}
func main() {

	validateArgs()
	dialer := &net.Dialer{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Timeout)*time.Second)
	conn, err := dialer.DialContext(ctx, "tcp", Host+":"+Port)
	if err != nil {
		log.Fatal("Failed to connect", err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		readSocket(ctx, cancel, conn)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		writeSocket(ctx, conn)
	}()
	wg.Wait()
	conn.Close()
}
