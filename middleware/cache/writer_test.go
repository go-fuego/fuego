package cache

import (
	"bytes"
	"net/http/httptest"
	"testing"
)

func TestWriter(t *testing.T) {
	w := httptest.NewRecorder()
	bytesBuffer := &bytes.Buffer{}
	m := &MultiHTTPWriter{
		ResponseWriter: w,
		cacheWriter:    bytesBuffer,
	}
	written := "hello world"

	n, err := m.Write([]byte(written))
	if err != nil {
		t.Fatal(err)
	}
	if len(written) != n {
		t.Errorf("Expected %d, got %d", len(written), n)
	}

	bytesWritten := m.cacheWriter.(*bytes.Buffer).String()
	if written != bytesWritten {
		t.Errorf("Expected %s, got %s", written, m.cacheWriter.(*bytes.Buffer).Bytes())
	}
}
