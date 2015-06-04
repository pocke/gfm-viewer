package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"

	"github.com/google/go-github/github"
)

type Storage struct {
	files map[string]string
	mu    *sync.RWMutex

	token   *Token
	watcher *Watcher
}

func NewStorage() *Storage {
	w, err := NewWatcher()
	if err != nil {
		panic(err)
	}

	s := &Storage{
		files:   make(map[string]string),
		token:   &Token{},
		mu:      &sync.RWMutex{},
		watcher: w,
	}

	go func() {
		ch := w.OnUpdate()
		for {
			fname := <-ch
			fmt.Println(fname)
		}
	}()

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
		err := s.watcher.AddFile(path)
		if err != nil {
			s.files[path] = err.Error()
			continue
		}
		s.AddFile(path)
	}
}

// without mutex
func (s *Storage) AddFile(path string) error {
	md, err := ioutil.ReadFile(path)
	if err != nil {
		s.files[path] = err.Error()
		return err
	}

	html, err := s.md2html(string(md))
	if err != nil {
		s.files[path] = html
		return err
	}
	html = s.insertCSS(html)
	s.files[path] = html
	return nil
}

func (s *Storage) UpdateAll() {
	s.AddFiles(s.Index())
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

func (_ *Storage) insertCSS(html string) string {
	tags := `<!DOCTYPE html>
<link rel="stylesheet" href="/css/github-markdown.css">
<div class="markdown-body">
<style>
.markdown-body { min-width: 200px; max-width: 790px; margin: 0 auto; padding: 30px; }
</style>
`
	tagEnd := `
</div>`
	return tags + html + tagEnd
}
