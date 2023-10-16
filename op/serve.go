package op

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func Run(port string) {
	http.ListenAndServe(port, nil)
}

type Controller[ReturnType any, Body any] func(c Ctx[Body]) (ReturnType, error)

func HttpHandler[ReturnType any, Body any](controller any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		f, ok := controller.(func(c Ctx[Body]) (ReturnType, error))
		if !ok {
			var c Controller[ReturnType, Body]
			slog.Info("Controller types not ok",
				"type", fmt.Sprintf("%T", controller),
				"should be", fmt.Sprintf("%T", c))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := Ctx[Body]{
			request: r,
		}

		ans, err := f(ctx)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
