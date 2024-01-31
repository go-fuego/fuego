package main

import "net/http"

type MyError struct {
	Err error // developer readable error message
}

func (e MyError) Status() int {
	return http.StatusTeapot
}
