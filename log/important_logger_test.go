package log_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/r3code/go-useful-snippets/log"
)

func fakeNow() time.Time {
	timeNow, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	return timeNow
}

func Test_Logger_Log(t *testing.T) {
	var buf bytes.Buffer
	logger := log.NewLogger(&buf, 1, false)

	if err := logger.Log("hello", "world"); err != nil {
		t.Fatal(err)
	}
	want := "hello=world\n"
	have := buf.String()
	if !strings.Contains(have, want) {
		t.Errorf("\nwant %#v\nhave %#v", want, have)
	}
	if !strings.Contains(have, "time=") {
		t.Errorf("No time= field printed")
	}
	if !strings.Contains(have, "caller=") {
		t.Errorf("No caller= field printed")
	}
}

type failingWriter struct {
	io.Writer
}

func (w *failingWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("Write failed")
}

func Test_Logger_Log_WriteError(t *testing.T) {
	var w failingWriter
	logger := log.NewLogger(&w, 3, false)

	var initPainc = func() {
		for i := 0; i <= 4; i++ {
			_ = logger.Log("hello", "world")
		}
	}
	// отложим проверку на возникновение паники, затем вызовем функцию в которой должна быть паника. Если паника не произошла - тест провалится
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Logger did not panic after 3 errors as it should")
		}
	}()

	initPainc()
}
