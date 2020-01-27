package main

import (
	"io"
	"strings"
	"testing"
)

// Verify that readFromIO catches EOF from reader and forward it to the ErrorChannel
func Test_readFromIO(t *testing.T) {

	testString := "Some test string"

	errCh := make(chan error)
	in := make(chan string)
	out := make(chan string)

	var err error

	go readFromIO(strings.NewReader(testString), in, out, errCh, true)

	for {
		select {
		// Take care about Error Channel, just to make shure that we get EOF and able to quit.
		case err = <-errCh:
			if err != io.EOF {
				t.Error("Expected", err, io.EOF)
			}
			return
		}
	}

}
