package main

import (
	"reflect"
	"sort"
	"sync"
	"testing"
)

type testWatcher struct {
	ch chan string
}

// Dummy methods
func (w *testWatcher) AddFile(_ string) error  { return nil }
func (w *testWatcher) Close() error            { return nil }
func (w *testWatcher) OnUpdate() <-chan string { return w.ch }

func (w *testWatcher) getCh() chan string { return w.ch }

type chGetter interface {
	getCh() chan string
}

var _ Watcher = &testWatcher{}
var _ chGetter = &testWatcher{}

func testStorage() *Storage {
	s := &Storage{
		files: make(map[string]file),
		token: &Token{},
		mu:    &sync.RWMutex{},
		watcher: &testWatcher{
			ch: make(chan string),
		},
		onUpdate: make(chan string),
	}
	go s.watch()

	return s
}

func TestStorageIndex(t *testing.T) {
	restore := helperSetEnv()
	defer restore()

	s := testStorage()
	files := []string{"hoge", "fuga", "poyo"}
	s.AddFiles(files)

	sort.Strings(files)
	if !reflect.DeepEqual(files, s.Index()) {
		t.Errorf("Expected: %v, but got %v", files, s.Index())
	}
}

func TestStorageGet(t *testing.T) {
	restore := helperSetEnv()
	defer restore()

	s := testStorage()
	path := "hogege"
	s.AddFiles([]string{path})

	_, exist := s.Get(path)
	if !exist {
		t.Error("File should be exist.")
	}
}

func TestStorageOnUpdate(t *testing.T) {
	restore := helperSetEnv()
	defer restore()

	s := testStorage()

	onUpdate := s.OnUpdate()
	w, _ := s.watcher.(chGetter)
	trigger := w.getCh()

	path := "poyoyoyoyo"

	trigger <- path

	got := <-onUpdate

	if got != path {
		t.Errorf("Expected: %s, but got %s", path, got)
	}
}
