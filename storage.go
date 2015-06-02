package main

import (
	"sort"
	"sync"
)

type Storage struct {
	files map[string]string
	mu    *sync.RWMutex

	token *Token
}

func NewStorage() *Storage {
	s := &Storage{
		files: make(map[string]string),
	}
	return s
}

func (s *Storage) Add(path, html string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.files[path] = html
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
