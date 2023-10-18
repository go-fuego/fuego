// Description: This file contains functions for sending JSON and XML responses.
package op

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Message    string `json:"error" xml:"Error"` // human readable error message
	StatusCode int    `json:"-" xml:"-"`         // http status code
}

var _ ErrorWithStatus = ErrorResponse{}

func (e ErrorResponse) Error() string {
	return e.Message
}

func (e ErrorResponse) Status() int {
	if e.StatusCode == 0 {
		return http.StatusInternalServerError
	}
	return e.StatusCode
}

// Send sends a string response.
func Send(w http.ResponseWriter, text string) {
	_, _ = w.Write([]byte(text))
}

// SendJSON sends a JSON response.
func SendJSON(w http.ResponseWriter, ans any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(ans)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Server Error"}`))
		return
	}
}

// SendJSONError sends a JSON error response.
// If the error implements ErrorWithStatus, the status code will be set.
func SendJSONError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	errResponse := ErrorResponse{
		Message: err.Error(),
	}

	status := http.StatusInternalServerError
	var errorStatus ErrorWithStatus
	if errors.As(err, &errorStatus) {
		status = errorStatus.Status()
	}

	w.WriteHeader(status)
	SendJSON(w, errResponse)
}

// SendXML sends a XML response.
func SendXML(w http.ResponseWriter, ans any) {
	w.Header().Set("Content-Type", "application/xml")
	err := xml.NewEncoder(w).Encode(ans)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Server Error"}`))
		return
	}
}

// SendXMLError sends a XML error response.
// If the error implements ErrorWithStatus, the status code will be set.
func SendXMLError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/xml")
	errResponse := ErrorResponse{
		Message: err.Error(),
	}

	status := http.StatusInternalServerError
	var errorStatus ErrorWithStatus
	if errors.As(err, &errorStatus) {
		status = errorStatus.Status()
	}

	w.WriteHeader(status)
	SendXML(w, errResponse)
}
