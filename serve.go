package fuego

import (
	"html/template"
	"log/slog"
	"net/http"
	"reflect"
	"time"
)

// Run starts the server.
func (s *Server) Run() {
	go s.generateOpenAPI()
	elapsed := time.Since(s.startTime)
	slog.Debug("Server started in "+elapsed.String(), "info", "time between since server creation (fuego.NewServer) and server startup (fuego.Run). Depending on your implementation, there might be things that do not depend on fuego slowing start time")
	slog.Info("Server running ✅ on http://localhost"+s.Server.Addr, "started in", elapsed.String())

	s.Server.Handler = s.Mux
	err := s.Server.ListenAndServe()
	if err != nil {
		slog.Error("Error running server", "error", err)
	}
}

type Controller[ReturnType any, Body any] func(c Context[Body]) (ReturnType, error)

// httpHandler converts a framework controller into a http.HandlerFunc.
func httpHandler[ReturnType any, Body any](s *Server, controller func(c Ctx[Body]) (ReturnType, error)) http.HandlerFunc {
	returnsHTML := reflect.TypeOf(controller).Out(0).Name() == "HTML"

	baseCtx := NewContext[Body](nil, nil, readOptions{
		DisallowUnknownFields: s.DisallowUnknownFields,
		MaxBodySize:           s.maxBodySize,
	})
	baseCtx.fs = s.fs
	if s.template != nil {
		baseCtx.templates = template.Must(s.template.Clone())
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := baseCtx.SafeShallowCopy()
		ctx.response = w
		ctx.request = r

		for _, param := range parsePathParams(r.URL.Path) {
			ctx.pathParams[param] = "coming in go1.22"
		}

		ans, err := controller(ctx)
		if err != nil {
			err = s.ErrorHandler(err)
			s.SerializeError(w, err)
			return
		}

		if reflect.TypeOf(ans) == nil {
			return
		}

		if returnsHTML {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, err = w.Write([]byte(any(ans).(HTML)))
			if err != nil {
				s.SerializeError(w, err)
			}
			return
		}

		ans, err = transformOut(r.Context(), ans)
		if err != nil {
			err = s.ErrorHandler(err)
			s.SerializeError(w, err)
			return
		}

		s.Serialize(w, ans)
	}
}
