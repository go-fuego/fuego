package fuego

import (
	"html/template"
	"log/slog"
	"net/http"
	"reflect"
	"time"
)

// Run starts the server.
// It is blocking.
// It returns an error if the server could not start (it could not bind to the port).
func (s *Server) Run() error {
	go s.generateOpenAPI()
	elapsed := time.Since(s.startTime)
	slog.Debug("Server started in "+elapsed.String(), "info", "time between since server creation (fuego.NewServer) and server startup (fuego.Run). Depending on your implementation, there might be things that do not depend on fuego slowing start time")
	slog.Info("Server running ✅ on http://localhost"+s.Server.Addr, "started in", elapsed.String())

	s.Server.Handler = s.Mux

	return s.Server.ListenAndServe()
}

// httpHandler converts a Fuego controller into a http.HandlerFunc.
func httpHandler[ReturnType any, Body any](s *Server, controller func(c Ctx[Body]) (ReturnType, error)) http.HandlerFunc {
	returnsHTML := reflect.TypeOf(controller).Out(0).Name() == "HTML"
	var r ReturnType
	_, returnsString := any(r).(*string)
	if !returnsString {
		_, returnsString = any(r).(string)
	}

	readOptions := readOptions{}
	if s != nil {
		readOptions.DisallowUnknownFields = s.DisallowUnknownFields
		readOptions.MaxBodySize = s.maxBodySize
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "Fuego")

		w.Header().Set("Trailer", "Server-Timing")
		timeCtxInit := time.Now()

		var templates *template.Template
		if s.template != nil {
			templates = template.Must(s.template.Clone())
		}
		ctx := &BaseContextWithBody[Body]{
			BaseContext: BaseContext{
				Req:         r,
				Res:         w,
				readOptions: readOptions,
				fs:          s.fs,
				templates:   templates,
			},
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

		if returnsString {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			stringToWrite, ok := any(ans).(string)
			if !ok {
				stringToWrite = *any(ans).(*string)
			}
			_, err = w.Write([]byte(stringToWrite))
			if err != nil {
				s.SerializeError(w, err)
			}
			w.Header().Add("Server-Timing", Timing{"write", time.Since(timeTransformOut), "transformOut"}.String())
			return
		}

		timeAfterTransformOut := time.Now()
		w.Header().Add("Server-Timing", Timing{"transformOut", timeAfterTransformOut.Sub(timeTransformOut), "transformOut"}.String())

		s.Serialize(w, ans)
		w.Header().Add("Server-Timing", Timing{"serialize", time.Since(timeAfterTransformOut), ""}.String())
	}
}
