package main

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
)

func Test_readerWriter(t *testing.T) {

	testCases := []struct{ description, in, out string }{
		{
			description: "Test with ordinary string",
			in:          "Test",
			out:         "Test\n",
		},
		{
			description: "Test with empty string",
			in:          "",
			out:         "\n",
		},
		{
			description: "Test with one space string",
			in:          " ",
			out:         " \n",
		},
	}

	ctx := context.Background()
	ctxWCancel, cancel := context.WithCancel(ctx)

	for _, tc := range testCases {
		buffer := &bytes.Buffer{}
		err := rw(ctxWCancel, cancel, strings.NewReader(tc.in), buffer)
		if err == io.EOF && buffer.String() != tc.out {
			t.Error("Error", buffer.String(), tc.in)
		}

	}

}
