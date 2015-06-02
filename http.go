package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/pocke/gfm-viewer/env"
	"github.com/yosssi/ace"
)

type Server struct {
	pages map[string]string
	mu    *sync.RWMutex
}

func NewServer() *Server {
	s := &Server{
		pages: make(map[string]string),
		mu:    &sync.RWMutex{},
	}

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if path == "/" || path == "/index.html" {
				s.indexHandler(w, r)
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

func (s *Server) Add(path, html string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pages[path] = html
}

// Update same as Add.
func (s *Server) Update(path, html string) {
	s.Add(path, html)
}

func (s *Server) Get(path string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	html, ok := s.pages[path]
	if ok {
		return html, ok
	} else {
		html, ok := s.pages["/"+path]
		return html, ok
	}
}

func (s *Server) Index() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]string, 0, len(s.pages))

	for path := range s.pages {
		res = append(res, path)
	}

	sort.Strings(res)
	return res
}

func (s *Server) ServeFile(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/files")
	html, ok := s.Get(path)
	if !ok {
		http.Error(w, fmt.Sprintf("%s page not found", path), http.StatusNotFound)
		return
	}
	w.Write([]byte(html))
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := ace.Load("assets/index", "", &ace.Options{
		DynamicReload: env.DEBUG,
		Asset:         Asset,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, s.Index())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
