package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/naoina/denco"
	"github.com/pocke/hlog"
	"github.com/skratchdot/open-golang/open"
	"github.com/yosssi/ace"
)

type Server struct {
	storage *Storage
}

func NewServer(port int) *Server {
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
			mux.GET("/ws", func(w http.ResponseWriter, r *http.Request, _ denco.Params) { wsm.ServeHTTP(w, r) }),
			mux.GET("/:type/:fname", s.serveAsset),
		})
		if err != nil {
			panic(err)
		}
		handler := f.ServeHTTP
		if DEBUG {
			handler = hlog.Wrap(f.ServeHTTP)
		}
		http.HandleFunc("/", handler)
		// TODO: port
		url := fmt.Sprintf("http://localhost:%d", port)
		fmt.Printf("Open: %s\n", url)
		open.Start(url)
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()

	return s
}

func (s *Server) ServeFile(w http.ResponseWriter, r *http.Request, p denco.Params) {
	path := p.Get("path")
	f, exist := s.storage.Get(path)
	if !exist {
		http.Error(w, fmt.Sprintf("%s page not found", path), http.StatusNotFound)
		return
	}
	if f.err != nil {
		http.Error(w, f.err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(f.html))
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
	tpl, err := ace.Load("assets/base", "assets/"+action, &ace.Options{
		DynamicReload: DEBUG,
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

func (_ *Server) serveAsset(w http.ResponseWriter, r *http.Request, p denco.Params) {
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
