package op

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

func (s *Server) Run() {
	slog.Info("Server running on http://localhost" + s.Config.Addr)
	_ = http.ListenAndServe(s.Config.Addr, s.mux)
}

type Controller[ReturnType any, Body any] func(c Context[Body]) (ReturnType, error)

type ErrorResponse struct {
	Error string `json:"error"` // human readable error message
}

// httpHandler converts a controller into a http.HandlerFunc.
func httpHandler[ReturnType any, Body any](controller func(c Ctx[Body]) (ReturnType, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context[Body]{
			request: r,
		}

		ans, err := controller(ctx)
		if err != nil {
			slog.Error("Error in controller", "err", err.Error())
			errResponse := ErrorResponse{
				Error: err.Error(),
			}

			status := http.StatusInternalServerError
			var errorStatus ErrorWithStatus
			if errors.As(err, &errorStatus) {
				status = errorStatus.Status()
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			err = json.NewEncoder(w).Encode(errResponse)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"Internal Server Error"}`))
				return
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(ans)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"Internal Server Error"}`))
			return
		}
	}
}
