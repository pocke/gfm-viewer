package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/naoina/denco"
	"github.com/pocke/gfm-viewer/env"
	"github.com/pocke/hlog"
	"github.com/yosssi/ace"
)

type Server struct {
	storage *Storage
}

func NewServer() *Server {
	s := &Server{
		storage: NewStorage(),
	}

	go func() {
		wsm := NewWSManager(s.storage.OnUpdate())

		mux := denco.NewMux()
		f, err := mux.Build([]denco.Handler{
			mux.GET("/", s.indexHandler),
			mux.POST("/auth", s.authHandler),
			mux.GET("/files/*path", s.ServeFile),
			mux.GET("/ws/*path", wsm.ServeWS),
			mux.GET("/:type/:fname", s.serveAsset),
		})
		if err != nil {
			panic(err)
		}
		handler := f.ServeHTTP
		if env.DEBUG {
			handler = hlog.Wrap(f.ServeHTTP)
		}
		http.HandleFunc("/", handler)
		// TODO: port
		http.ListenAndServe(":1124", nil)
	}()

	return s
}

func (s *Server) ServeFile(w http.ResponseWriter, r *http.Request, p denco.Params) {
	path := p.Get("path")
	html, ok := s.storage.Get(path)
	if !ok {
		http.Error(w, fmt.Sprintf("%s page not found", path), http.StatusNotFound)
		return
	}
	w.Write([]byte(html))
}

func (s *Server) authHandler(w http.ResponseWriter, r *http.Request, _ denco.Params) {
	r.ParseForm()
	v := r.PostForm
	user := v.Get("username")
	pass := v.Get("password")

	err := s.storage.token.Init(user, pass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.storage.AddAll()
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request, _ denco.Params) {
	if s.storage.token.hasToken() {
		loadAce(w, "index", s.storage.Index())
	} else {
		loadAce(w, "before_auth", nil)
	}
}

func loadAce(w http.ResponseWriter, action string, data interface{}) {
	tpl, err := ace.Load("assets/"+action, "", &ace.Options{
		DynamicReload: env.DEBUG,
		Asset:         Asset,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) serveAsset(w http.ResponseWriter, r *http.Request, p denco.Params) {
	t := p.Get("type")
	fname := p.Get("fname")
	file, err := Asset(path.Join("assets", t, fname))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var contentType string
	switch t {
	case "js":
		contentType = "application/javascript"
	case "css":
		contentType = "text/css"
	}
	w.Header().Set("Content-Type", contentType)
	w.Write(file)
}
