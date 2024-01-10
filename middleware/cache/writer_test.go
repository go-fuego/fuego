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

func TestWriter_Unwrap(t *testing.T) {
	w := httptest.NewRecorder()
	bytesBuffer := &bytes.Buffer{}
	m := &MultiHTTPWriter{
		ResponseWriter: w,
		cacheWriter:    bytesBuffer,
	}
	if m.Unwrap() != w {
		t.Errorf("Expected %T, got %T", w, m.Unwrap())
	}
}

func TestWriter_WriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	bytesBuffer := &bytes.Buffer{}
	m := &MultiHTTPWriter{
		ResponseWriter: w,
		cacheWriter:    bytesBuffer,
	}
	m.WriteHeader(204)
	if m.status != 204 {
		t.Errorf("Expected %d, got %d", 204, m.status)
	}
}
