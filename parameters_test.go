package main

import (
	"testing"
)

func TestValidate(t *testing.T) {

	testCases := []struct {
		arg []string
		err error
	}{
		{
			arg: []string{"--timeout=10s", "localhost", "8081"},
			err: nil,
		},
		{
			arg: []string{"--timeout=10s", "localhost", "abc123"},
			err: ErrPortMustBeANumber,
		},
		{
			arg: []string{"--timeout=10s"},
			err: ErrHostnPortMustBeDefined,
		},
		{
			arg: []string{"localhost"},
			err: ErrHostnPortMustBeDefined,
		},
		{
			arg: []string{"123"},
			err: ErrHostnPortMustBeDefined,
		},
		{
			arg: []string{"--timeout=10s", "123"},
			err: ErrHostnPortMustBeDefined,
		},
		{
			arg: []string{"--timeout=10", "localhost", "123"},
			err: ErrInvalidTimeFormat,
		},
	}

	for _, tc := range testCases {

		if err := validateArgs(tc.arg); err != tc.err {
			t.Error("Error", tc.arg, tc, err)
		}
	}

}
