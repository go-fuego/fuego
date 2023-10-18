package op

import (
	"net/http"

	"log/slog"
)

func (s *Server) Run() {
	s.GenerateOpenAPI()
	slog.Info("Server running on http://localhost" + s.Addr)
	_ = http.ListenAndServe(s.Addr, s.mux)
}

type Controller[ReturnType any, Body any] func(c Context[Body]) (ReturnType, error)

// httpHandler converts a controller into a http.HandlerFunc.
func httpHandler[ReturnType any, Body any](s *Server, controller func(c Ctx[Body]) (ReturnType, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context[Body]{
			request: r,
			readOptions: readOptions{
				DisallowUnknownFields: s.DisallowUnknownFields,
				MaxBodySize:           s.maxBodySize,
			},
		}

		ans, err := controller(ctx)
		if err != nil {
			s.SerializeError(w, err)
			return
		}

		s.Serialize(w, ans)
	}
}
