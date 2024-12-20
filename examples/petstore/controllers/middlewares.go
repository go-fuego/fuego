package controller

import "net/http"

func dummyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do something before
		next.ServeHTTP(w, r)
		// do something after
	})
}
