package op

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

func Run(port string) {
	http.ListenAndServe(config.Addr, defaultMux)
}

type Controller[ReturnType any, Body any] func(c Ctx[Body]) (ReturnType, error)

type ErrorResponse struct {
	Error string `json:"error"` // human readable error message
}

// httpHandler converts a controller into a http.HandlerFunc.
func httpHandler[ReturnType any, Body any](controller func(c Ctx[Body]) (ReturnType, error)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := Ctx[Body]{
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
			json.NewEncoder(w).Encode(errResponse)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(ans)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
