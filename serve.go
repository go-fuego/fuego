package op

import (
	"log/slog"
	"net/http"
	"time"
)

func (s *Server) Run() {
	go s.GenerateOpenAPI()
	elapsed := time.Since(s.startTime)
	slog.Debug("Server started in "+elapsed.String(), "info", "time between since server creation (op.NewServer) and server startup (op.Run). Depending on your implementation, there might be things that do not depend on op slowing start time")
	slog.Info("Server running ✅ on http://localhost"+s.Addr, "started in", elapsed.String())

	server := &http.Server{
		Addr:              s.Addr,
		Handler:           s.mux,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("Error running server", "error", err)
	}
}

type Controller[ReturnType any, Body any] func(c Context[Body]) (ReturnType, error)

// httpHandler converts a controller into a http.HandlerFunc.
func httpHandler[ReturnType any, Body any](s *Server, controller func(c Ctx[Body]) (ReturnType, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext[Body](r, readOptions{
			DisallowUnknownFields: s.DisallowUnknownFields,
			MaxBodySize:           s.maxBodySize,
		})

		ans, err := controller(ctx)
		if err != nil {
			s.SerializeError(w, err)
			return
		}

		ans, err = transformOut(r.Context(), ans)
		if err != nil {
			s.SerializeError(w, err)
			return
		}

		s.Serialize(w, ans)
	}
}
