package main

import "github.com/go-fsnotify/fsnotify"

type Watcher struct {
	w        *fsnotify.Watcher
	onUpdate chan string
}

func NewWatcher() (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	res := &Watcher{
		w:        w,
		onUpdate: make(chan string),
	}

	go func() {
		for {
			select {
			case ev := <-w.Events:
				switch ev.Op {
				case fsnotify.Rename:
					continue
				case fsnotify.Remove:
					w.Add(ev.Name)
					continue
				}
				res.onUpdate <- ev.Name
			case err := <-w.Errors:
				panic(err)
			}
		}
	}()

	return res, nil
}

func (w *Watcher) AddFile(path string) error {
	return w.w.Add(path)
}

// TODO: フルパスが返ってきたりしたらアレっぽい
func (w *Watcher) OnUpdate() <-chan string { return w.onUpdate }
func (w *Watcher) Close() error            { return w.w.Close() }
