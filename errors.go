package op

import "net/http"

// ErrorWithStatus is an interface that can be implemented by an error to provide
// additional information about the error.
type ErrorWithStatus interface {
	Status() int
}

// ErrorWithInfo is an interface that can be implemented by an error to provide
// additional information about the error.
type ErrorWithInfo interface {
	Info() map[string]any
}

// ErrorResponse is the error response used by the serialization part of the framework.
type ErrorResponse struct {
	Message    string         `json:"error" xml:"Error"`                   // human readable error message
	StatusCode int            `json:"-" xml:"-"`                           // http status code
	MoreInfo   map[string]any `json:"info,omitempty" xml:"Info,omitempty"` // additional info
}

func (e ErrorResponse) Error() string {
	return e.Message
}

var _ ErrorWithStatus = ErrorResponse{}

func (e ErrorResponse) Status() int {
	if e.StatusCode == 0 {
		return http.StatusInternalServerError
	}
	return e.StatusCode
}

var _ ErrorWithInfo = ErrorResponse{}

func (e ErrorResponse) Info() map[string]any {
	return e.MoreInfo
}
