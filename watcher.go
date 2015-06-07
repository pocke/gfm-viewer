package main

import (
	"sync"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

type Watcher struct {
	w        *fsnotify.Watcher
	onUpdate chan string
	buf      chan string
}

func NewWatcher() (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	res := &Watcher{
		w:        w,
		onUpdate: make(chan string),
		buf:      make(chan string, 3),
	}

	go func() {
		for {
			select {
			case ev := <-w.Events:
				if ev.Op == fsnotify.Remove {
					w.Add(ev.Name)
				}
				res.buf <- ev.Name
			case err := <-w.Errors:
				if err != nil {
					panic(err)
				}
			}
		}
	}()

	go res.watchBuffer()

	return res, nil
}

func (w *Watcher) AddFile(path string) error {
	return w.w.Add(path)
}

// OnUpdate returns channel that notify on file update.
func (w *Watcher) OnUpdate() <-chan string { return w.onUpdate }

func (w *Watcher) Close() error { return w.w.Close() }

// Vim notifies three times.
// Ref: Japanese blog: http://qiita.com/ma2saka/items/d30e48b4c72f1f5f4873
// So, watchBuffer packs notifies around 50 Milli Second.
func (w *Watcher) watchBuffer() {
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
