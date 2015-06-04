package main

import "gopkg.in/fsnotify.v1"

type Watcher struct {
	w        *fsnotify.Watcher
	onWrite  chan string
	onRemove chan string
}

func NewWatcher() (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	res := &Watcher{
		w: w,
	}

	go func() {
		for {
			select {
			case ev := <-w.Events:
				switch ev.Op {
				case fsnotify.Write:
					res.onWrite <- ev.Name
				case fsnotify.Remove:
					res.onRemove <- ev.Name
				}
			case err := <-w.Errors:
				panic(err)
			}
		}
	}()

	return res, nil
}

func (w *Watcher) AddFiles(paths []string) error {
	for _, v := range paths {
		err := w.w.Add(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO: フルパスが返ってきたりしたらアレっぽい
func (w *Watcher) OnWrite() <-chan string  { return w.onWrite }
func (w *Watcher) OnRemove() <-chan string { return w.onRemove }
func (w *Watcher) Close() error            { return w.w.Close() }
