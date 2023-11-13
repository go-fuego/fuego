package cache

import (
	"io"
	"net/http"
)

// MultiHTTPWriter is a http.ResponseWriter that writes the response to multiple writers
type MultiHTTPWriter struct {
	http.ResponseWriter
	status      int       // status is the status code that will be written to the response
	cacheWriter io.Writer // cacheWriter is the writer that will be used to cache the response
}

var _ http.ResponseWriter = &MultiHTTPWriter{}

func (m *MultiHTTPWriter) Write(p []byte) (int, error) {
	multiWriter := io.MultiWriter(m.ResponseWriter, m.cacheWriter)
	return multiWriter.Write(p)
}

func (m *MultiHTTPWriter) Unwrap() http.ResponseWriter {
	return m.ResponseWriter
}
func (m *MultiHTTPWriter) WriteHeader(statusCode int) {
	m.status = statusCode
	m.ResponseWriter.WriteHeader(statusCode)
}
