package main

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
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
	templates map[string]*template.Template
}

func (t *Renderer) AddPage(templ string) error {
	tt, err := template.ParseGlob("views/" + templ)
	if err != nil {
		panic(err)
	}
	tt, err = tt.New(templ).ParseFiles("views/htmx.layout.html")
	if err != nil {
		panic(err)
	}
	t.templates[templ] = tt
	fullTempl := "full-" + templ
	ft, err := tt.New(fullTempl).ParseFiles("views/application.layout.html")
	if err != nil {
		panic(err)
	}

	t.templates[fullTempl] = ft
	return nil
}

func (ren *Renderer) Render(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	// if this view doesn't exist try adding it first
	if _, ok := ren.templates[name]; !ok {
		ren.AddPage(name)
	}
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

var k = koanf.New(".")

func LoadConfigs() error {
	if err := k.Load(file.Provider("./config.yaml"), yaml.Parser()); err != nil {
		return err
	}

	k.Load(env.Provider("CONTACTS_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "CONTACTS_")), "_", ".", -1)
	}), nil)

	return nil
}

func main() {
	if err := LoadConfigs(); err != nil {
		panic(err)
	}

	l := mustCreateLogger(k.String("contacts.log_level"), k.String("contacts.log_encoding"))
	renderer := &Renderer{l: l, templates: make(map[string]*template.Template)}
	fmt.Println(renderer)
	r := chi.NewRouter()
	r.Use(middleware.Compress(5))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/contacts", http.StatusPermanentRedirect)
		//		renderer.Render(w, r, "contacts/index.html", nil)
	})

	r.Get("/contacts", func(w http.ResponseWriter, r *http.Request) {
		renderer.Render(w, r, "contacts/index.html", nil)
	})

	srv := &http.Server{
		Addr:    net.JoinHostPort(k.String("contacts.host"), k.String("contacts.port")),
		Handler: r,
	}

	l.Info("starting server", zap.String("host", k.String("contacts.host")), zap.String("port", k.String("contacts.port")))

	err := srv.ListenAndServe()
	if err != nil {
		l.Error("error listening", zap.Error(err))
	}
}
