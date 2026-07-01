package peekio_test

import (
	"bufio"
	"io"
	"strings"
	"testing"

	"github.com/alexatcanva/peekio"
)

func TestPeekReader_Read(t *testing.T) {
	testCases := []struct {
		desc        string
		haveContent string
		haveReadLen int
		wantBytes   []byte
		wantN       int
		wantError   error
	}{
		{
			desc:        "reads first 5 bytes",
			haveContent: "Hello, World!",
			haveReadLen: 5,
			wantBytes:   []byte("Hello"),
			wantN:       5,
			wantError:   nil,
		},
		{
			desc:        "empty reader returns EOF",
			haveContent: "",
			haveReadLen: 5,
			wantBytes:   []byte{},
			wantN:       0,
			wantError:   io.EOF,
		},
		{
			desc:        "read past end returns available bytes without error",
			haveContent: "Hi",
			haveReadLen: 10,
			wantBytes:   []byte("Hi"),
			wantN:       2,
			wantError:   nil,
		},
		{
			desc:        "read exactly content length",
			haveContent: "Hello",
			haveReadLen: 5,
			wantBytes:   []byte("Hello"),
			wantN:       5,
			wantError:   nil,
		},
		{
			desc:        "single byte read",
			haveContent: "Hello",
			haveReadLen: 1,
			wantBytes:   []byte("H"),
			wantN:       1,
			wantError:   nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			underlying := bufio.NewReader(strings.NewReader(tC.haveContent))
			peeker := peekio.NewPeekReader(underlying)

			buf := make([]byte, tC.haveReadLen)
			n, gotErr := peeker.Read(buf)

			if n != tC.wantN || string(buf[:n]) != string(tC.wantBytes) || gotErr != tC.wantError {
				t.Errorf("Read() = (%q, %v); want (%q, %v)", buf[:n], gotErr, tC.wantBytes, tC.wantError)
			}

			// Peek must not advance the underlying reader's position.
			afterBuf := make([]byte, len(tC.haveContent))
			underlying.Read(afterBuf)
			if string(afterBuf) != tC.haveContent {
				t.Errorf("underlying reader advanced after Read(): got %q, want %q", afterBuf, tC.haveContent)
			}
		})
	}
}

// TestPeekReader_SequentialReads verifies that successive reads advance through
// the content and that EOF is only returned once the content is exhausted.
func TestPeekReader_SequentialReads(t *testing.T) {
	underlying := bufio.NewReader(strings.NewReader("Hello, World!"))
	peeker := peekio.NewPeekReader(underlying)

	reads := []struct {
		wantBytes string
		wantError error
	}{
		{"Hello", nil},
		{", Wor", nil},
		{"ld!", nil}, // only 3 bytes remain, but read buf is 5
		{"", io.EOF}, // content exhausted
	}

	for i, r := range reads {
		buf := make([]byte, 5)
		n, err := peeker.Read(buf)
		if string(buf[:n]) != r.wantBytes || err != r.wantError {
			t.Errorf("Read() #%d = (%q, %v); want (%q, %v)", i+1, buf[:n], err, r.wantBytes, r.wantError)
		}
	}
}
