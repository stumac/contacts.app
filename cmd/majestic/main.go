package main

import (
	"html/template"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// make proper zap config based on level
func mustCreateLogger(level, encoding string) *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.Encoding = encoding
	zapLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		panic(err)
	}
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return l
}

type Renderer struct {
	l         *zap.Logger
	templates map[string]template.Template
}

func (ren *Renderer) Render(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	if r.Header.Get("HX-Request") == "true" {
		templ, ok := ren.templates[name]
		if !ok {
			ren.l.Error("html template not found!", zap.Bool("hxRequest", true), zap.String("viewName", name))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("something went wrong"))
			return
		}

		err := templ.ExecuteTemplate(w, "htmx.layout.html", data)
		if err != nil {
			ren.l.Error("error executing template", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("something went wrong"))
		}

	} else {
		templ, ok := ren.templates["full-"+name]
		if !ok {
			ren.l.Error("html template not found!", zap.Bool("hxRequest", true), zap.String("viewName", name))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("something went wrong"))
			return
		}

		err := templ.ExecuteTemplate(w, "application.layout.html", data)
		if err != nil {
			ren.l.Error("error executing template", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("something went wrong"))
		}
	}
}

func main() {
	l := mustCreateLogger(os.Getenv("LOG_LEVEL"), os.Getenv("LOG_ENCODING"))
	renderer := &Renderer{l: l, templates: make(map[string]template.Template)}
	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		renderer.Render(w, r, "testy", nil)
		w.Write([]byte("I hear xhtml is pretty based"))
	})

	srv := &http.Server{
		Addr:    net.JoinHostPort(os.Getenv("HOST"), os.Getenv("PORT")),
		Handler: r,
	}

	l.Info("starting server", zap.String("host", os.Getenv("HOST")), zap.String("port", os.Getenv("PORT")))

	err := srv.ListenAndServe()
	if err != nil {
		l.Error("error listening", zap.Error(err))
	}
}
