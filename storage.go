package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"

	"github.com/google/go-github/github"
)

type file struct {
	html string
	err  error
}

// Storage parses markdown. And save this. HTTP Server read parsed markdown from Storage.
type Storage struct {
	files map[string]file
	mu    *sync.RWMutex

	token   *Token
	watcher *Watcher
	// onUpdate notify when file is updated.
	onUpdate chan string
}

// NewStorage creates a new Storage.
// And watch changing file. When notify changing file, parse markdown and notify by 'onUpdate' channel.
func NewStorage() *Storage {
	w, err := NewWatcher()
	if err != nil {
		panic(err)
	}

	s := &Storage{
		files:    make(map[string]file),
		token:    &Token{},
		mu:       &sync.RWMutex{},
		watcher:  w,
		onUpdate: make(chan string),
	}

	go func() {
		ch := w.OnUpdate()
		for {
			fname := <-ch
			s.UpdateFile(fname)
		}
	}()

	return s
}

// AddFiles parses and saves parsed markdowns.
func (s *Storage) AddFiles(paths []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.token.hasToken() {
		for _, path := range paths {
			s.files[path] = file{
				err: errors.New("GitHub API token doesn't exist."),
			}
		}
		return
	}

	for _, path := range paths {
		err := s.watcher.AddFile(path)
		if err != nil {
			s.files[path] = file{err: err}
			continue
		}
		s.AddFile(path)
	}
}

// without mutex
func (s *Storage) AddFile(path string) {
	md, err := ioutil.ReadFile(path)
	if err != nil {
		s.files[path] = file{err: err}
		return
	}

	html, err := s.md2html(string(md))
	Log("Markdown parse request done for %s", path)
	if err != nil {
		s.files[path] = file{err: err}
		return
	}
	s.files[path] = file{html: html}
	return
}

// UpdateFile update saved file, and notify update.
func (s *Storage) UpdateFile(path string) {
	s.AddFile(path)
	s.onUpdate <- path
}

func (s *Storage) AddAll() {
	s.AddFiles(s.Index())
}

func (s *Storage) Get(path string) (file, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	f, exist := s.files[path]
	if exist {
		return f, exist
	} else {
		f, exist := s.files["/"+path]
		return f, exist
	}
}

// Index returns path list of saved files.
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

// md2html parses markdown.
func (s *Storage) md2html(md string) (string, error) {
	client := github.NewClient(&http.Client{
		Transport: s.token,
	})
	html, _, err := client.Markdown(md, nil)
	return html, err
}

func (s *Storage) OnUpdate() <-chan string {
	return s.onUpdate
}
