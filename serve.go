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
		w.Header().Set("Trailer", "Server-Timing")
		timeCtxInit := time.Now()

		ctx := baseCtx.SafeShallowCopy()
		ctx.response = w
		ctx.request = r

		for _, param := range parsePathParams(r.URL.Path) {
			ctx.pathParams[param] = "coming in go1.22"
		}

		timeController := time.Now()
		w.Header().Set("Server-Timing", Timing{"fuegoReqInit", timeController.Sub(timeCtxInit), ""}.String())

		ans, err := controller(ctx)
		if err != nil {
			err = s.ErrorHandler(err)
			s.SerializeError(w, err)
			return
		}
		timeAfterController := time.Now()
		w.Header().Add("Server-Timing", Timing{"controller", timeAfterController.Sub(timeController), ""}.String())

		if reflect.TypeOf(ans) == nil {
			return
		}

		ctxRenderer, ok := any(ans).(CtxRenderer)
		if ok {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			err = ctxRenderer.Render(r.Context(), w)
			if err != nil {
				err = s.ErrorHandler(err)
				s.SerializeError(w, err)
			}
			w.Header().Set("Server-Timing", Timing{"render", time.Since(timeAfterController), ""}.String())
			return
		}

		renderer, ok := any(ans).(Renderer)
		if ok {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			err = renderer.Render(w)
			if err != nil {
				err = s.ErrorHandler(err)
				s.SerializeError(w, err)
			}
			w.Header().Add("Server-Timing", Timing{"render", time.Since(timeAfterController), ""}.String())
			return
		}

		if returnsHTML {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, err = w.Write([]byte(any(ans).(HTML)))
			if err != nil {
				s.SerializeError(w, err)
			}
			w.Header().Add("Server-Timing", Timing{"render", time.Since(timeAfterController), ""}.String())
			return
		}

		timeTransformOut := time.Now()
		ans, err = transformOut(r.Context(), ans)
		if err != nil {
			err = s.ErrorHandler(err)
			s.SerializeError(w, err)
			return
		}
		timeAfterTransformOut := time.Now()
		w.Header().Add("Server-Timing", Timing{"transformOut", timeAfterTransformOut.Sub(timeTransformOut), "transformOut"}.String())

		s.Serialize(w, ans)
		w.Header().Add("Server-Timing", Timing{"serialize", time.Since(timeAfterTransformOut), ""}.String())
	}
}
