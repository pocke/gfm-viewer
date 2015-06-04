package main

import (
	"io/ioutil"
	"net/http"
	"sort"
	"sync"

	"github.com/google/go-github/github"
)

type Storage struct {
	files map[string]string
	mu    *sync.RWMutex

	token *Token
}

func NewStorage() *Storage {
	s := &Storage{
		files: make(map[string]string),
		token: &Token{},
		mu:    &sync.RWMutex{},
	}
	return s
}

func (s *Storage) AddFiles(paths []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.token.hasToken() {
		for _, path := range paths {
			s.files[path] = ""
		}
		return
	}

	for _, path := range paths {
		md, err := ioutil.ReadFile(path)
		if err != nil {
			s.files[path] = err.Error()
			continue
		}

		html, err := s.md2html(string(md))
		if err != nil {
			s.files[path] = html
			continue
		}
		s.files[path] = html
	}
}

func (s *Storage) Get(path string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	html, ok := s.files[path]
	if ok {
		return html, ok
	} else {
		html, ok := s.files["/"+path]
		return html, ok
	}
}

func (s *Storage) Index() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]string, 0, len(s.files))

	for path := range s.files {
		res = append(res, path)
	}

	sort.Strings(res)
	return res
}

func (s *Storage) md2html(md string) (string, error) {
	client := github.NewClient(&http.Client{
		Transport: s.token,
	})
	html, _, err := client.Markdown(md, nil)
	return html, err
}
