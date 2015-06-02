package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pocke/gfm-viewer/env"
	"github.com/yosssi/ace"
)

type Server struct {
	storage *Storage
}

func NewServer() *Server {
	s := &Server{
		storage: &Storage{},
	}

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if path == "/" || path == "/index.html" {
				s.beforeAuthHandler(w, r)
			} else if path == "/auth" {
				s.authHandler(w, r)
			} else if strings.HasPrefix(path, "/files") {
				s.ServeFile(w, r)
			} else {
				http.Error(w, "404 Not Found", http.StatusNotFound)
				return
			}
		})
		// TODO: port
		http.ListenAndServe(":1124", nil)
	}()

	return s
}

func (s *Server) ServeFile(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/files")
	html, ok := s.storage.Get(path)
	if !ok {
		http.Error(w, fmt.Sprintf("%s page not found", path), http.StatusNotFound)
		return
	}
	w.Write([]byte(html))
}

func (s *Server) authHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	v := r.PostForm
	user := v.Get("username")
	pass := v.Get("password")

	token, err := NewToken(user, pass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(token.Token))
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	loadAce(w, "index", s.storage.Index())
}

func (s *Server) beforeAuthHandler(w http.ResponseWriter, r *http.Request) {
	loadAce(w, "before_auth", nil)
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
