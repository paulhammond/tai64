// Copyright 2014 Paul Hammond.
// This software is licensed under the MIT license, see LICENSE.txt for details.

package tai64

import (
	"fmt"
	"testing"
	"time"
)

func ExampleParseTai64() {
	t, err := ParseTai64n("@4000000037c219bf2ef02e94")
	if err == nil {
		fmt.Println(t.UTC())
	}
	// Output:
	// 1999-08-24 04:03:43.7874925 +0000 UTC
}

var tai64nTests = []struct {
	hex   string
	bytes []byte
	time  string
}{
	// from `man 8 tai64nlocal`, converted to UTC
	{"@4000000037c219bf2ef02e94", []byte{0x40, 0x00, 0x00, 0x00, 0x37, 0xc2, 0x19, 0xbf, 0x2e, 0xf0, 0x2e, 0x94}, "1999-08-24T04:03:43.7874925Z"},
	// `echo @4000000052c65e550cd675fc | TZ=:/usr/share/zoneinfo/right/Etc/UTC tai64nlocal`
	{"@4000000052c65e550cd675fc", []byte{0x40, 0x00, 0x00, 0x00, 0x52, 0xc6, 0x5e, 0x55, 0x0c, 0xd6, 0x75, 0xfc}, "2014-01-03T06:52:34.2153815Z"},
	// the golang date, converted using http://www.tai64.com/
	{"@4000000043b9410600000000", []byte{0x40, 0x00, 0x00, 0x00, 0x43, 0xb9, 0x41, 0x06, 0x00, 0x00, 0x00, 0x00}, "2006-01-02T15:04:05Z"},

	// from http://cr.yp.to/libtai/tai64.html
	// the first second of 1970 TAI (which is 10 seconds earlier in UTC)
	{"@400000000000000000000000", []byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, "1969-12-31T23:59:50Z"},
	// one second later
	{"@400000000000000100000000", []byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}, "1969-12-31T23:59:51Z"},
	// 10 seconds later (the first second of 1970 UTC)
	{"@400000000000000A00000000", []byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0A, 0x00, 0x00, 0x00, 0x00}, "1970-01-01T00:00:00Z"},
	// the last second of 1969 TAI (which is 10 seconds earlier in UTC)
	{"@3FFFFFFFFFFFFFFF00000000", []byte{0x3F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00}, "1969-12-31T23:59:49Z"},
	// 1992-06-02 08:07:09 TAI
	{"@400000002a2b2c2d00000000", []byte{0x40, 0x00, 0x00, 0x00, 0x2a, 0x2b, 0x2c, 0x2d, 0x00, 0x00, 0x00, 0x00}, "1992-06-02T08:06:43Z"},
}

var tai64Tests = []struct {
	hex   string
	bytes []byte
	time  string
}{
	{"@4000000037c219bf", []byte{0x40, 0x00, 0x00, 0x00, 0x37, 0xc2, 0x19, 0xbf}, "1999-08-24T04:03:43Z"},
	{"@4000000052c65e55", []byte{0x40, 0x00, 0x00, 0x00, 0x52, 0xc6, 0x5e, 0x55}, "2014-01-03T06:52:34Z"},
	{"@4000000043b94106", []byte{0x40, 0x00, 0x00, 0x00, 0x43, 0xb9, 0x41, 0x06}, "2006-01-02T15:04:05Z"},
	{"@4000000000000000", []byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, "1969-12-31T23:59:50Z"},
	{"@4000000000000001", []byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, "1969-12-31T23:59:51Z"},
	{"@400000000000000A", []byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0A}, "1970-01-01T00:00:00Z"},
	{"@3FFFFFFFFFFFFFFF", []byte{0x3F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, "1969-12-31T23:59:49Z"},
	{"@400000002a2b2c2d", []byte{0x40, 0x00, 0x00, 0x00, 0x2a, 0x2b, 0x2c, 0x2d}, "1992-06-02T08:06:43Z"},
}

func TestParseTai64n(t *testing.T) {
	for _, test := range tai64nTests {
		result, err := ParseTai64n(test.hex)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if out := result.UTC().Format(time.RFC3339Nano); out != test.time {
			t.Errorf("got %v, expected %v", out, test.time)
		}
	}

	bad := []string{
		// no @
		"4000000037c219bf2ef02e94",
		"4000000037c219bf2ef02e941",
		// too short
		"@4000000037c219bf2ef02e9",
		// too long
		"@4000000037c219bf2ef02e941",
		// too big a number
		"@f000000037c219bf2ef02e94",
		// not hex
		"@G00000000000000000000000",
	}
	for _, test := range bad {
		result, err := ParseTai64n(test)
		if err != parseError {
			t.Errorf("expected %v, got %v", parseError, err)
		}
		if !result.IsZero() {
			t.Errorf("expected zero time, got %v", result)
		}
	}
}

func TestDecodeTai64n(t *testing.T) {
	for _, test := range tai64nTests {
		result, err := DecodeTai64n(test.bytes)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if out := result.UTC().Format(time.RFC3339Nano); out != test.time {
			t.Errorf("got %v, expected %v", out, test.time)
		}
	}
	bad := [][]byte{
		// too long
		[]byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		// too short
		[]byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		// too big a number
		[]byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
	for _, test := range bad {
		result, err := DecodeTai64n(test)
		if err != parseError {
			t.Errorf("expected %v, got %v", parseError, err)
		}
		if !result.IsZero() {
			t.Errorf("expected zero time, got %v", result)
		}
	}
}

func TestParseTai64(t *testing.T) {
	for _, test := range tai64Tests {
		result, err := ParseTai64(test.hex)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if out := result.UTC().Format(time.RFC3339); out != test.time {
			t.Errorf("got %v, expected %v", out, test.time)
		}
	}

	bad := []string{
		// no @
		"4000000037c219bf",
		"4000000037c219bf1",
		// too short
		"@4000000037c219b",
		// too long
		"@4000000037c219bf1",
		// too big a number
		"@f000000037c219bf",
		// not hex
		"@G000000000000000",
	}
	for _, test := range bad {
		result, err := ParseTai64n(test)
		if err != parseError {
			t.Errorf("expected %v, got %v", parseError, err)
		}
		if !result.IsZero() {
			t.Errorf("expected zero time, got %v", result)
		}
	}
}

func TestDecodeTai64(t *testing.T) {
	for _, test := range tai64Tests {
		result, err := DecodeTai64(test.bytes)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if out := result.UTC().Format(time.RFC3339); out != test.time {
			t.Errorf("got %v, expected %v", out, test.time)
		}
	}
	bad := [][]byte{
		// too long
		[]byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		// too short
		[]byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		// too big a number
		[]byte{0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
	for _, test := range bad {
		result, err := DecodeTai64n(test)
		if err != parseError {
			t.Errorf("expected %v, got %v", parseError, err)
		}
		if !result.IsZero() {
			t.Errorf("expected zero time, got %v", result)
		}
	}
}
