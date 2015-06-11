package main

import (
	"sync"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

type watcher struct {
	w        *fsnotify.Watcher
	onUpdate chan string
	buf      chan string
}

type Watcher interface {
	AddFile(string) error
	OnUpdate() <-chan string
	Close() error
}

func NewWatcher() (Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	res := &watcher{
		w:        w,
		onUpdate: make(chan string),
		buf:      make(chan string, 3),
	}

	go res.watchFS()
	go res.watchBuffer()

	return res, nil
}

func (w *watcher) AddFile(path string) error {
	return w.w.Add(path)
}

// OnUpdate returns channel that notify on file update.
func (w *watcher) OnUpdate() <-chan string { return w.onUpdate }

func (w *watcher) Close() error { return w.w.Close() }

// Vim notifies three times.
// Ref: Japanese blog: http://qiita.com/ma2saka/items/d30e48b4c72f1f5f4873
// So, watchBuffer packs notifies around 50 Milli Second.
func (w *watcher) watchBuffer() {
	mu := &sync.Mutex{}
	flags := make(map[string]bool)

	for {
		name := <-w.buf

		mu.Lock()
		if flags[name] {
			mu.Unlock()
			continue
		}

		flags[name] = true
		mu.Unlock()

		go func(n string) {
			<-time.After(50 * time.Millisecond)
			mu.Lock()
			defer mu.Unlock()
			flags[n] = false
			w.onUpdate <- n
		}(name)
	}
}

func (w *watcher) watchFS() {
	for {
		select {
		case ev := <-w.w.Events:
			if ev.Op == fsnotify.Remove {
				w.w.Add(ev.Name)
			}
			w.buf <- ev.Name
		case err := <-w.w.Errors:
			if err != nil {
				panic(err)
			}
		}
	}
}
