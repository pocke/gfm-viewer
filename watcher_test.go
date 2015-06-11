package main

import (
	"io/ioutil"
	"os"
	"testing"
)

var _ Watcher = &watcher{}

func TestWatcher(t *testing.T) {
	w, err := NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	files := make([]*os.File, 0, 2)
	for i := 0; i < 2; i++ {
		f, err := ioutil.TempFile("", "gfm-viewer-test")
		if err != nil {
			t.Fatal(err)
		}
		fname := f.Name()
		defer os.Remove(fname)

		err = w.AddFile(fname)
		if err != nil {
			t.Fatal(err)
		}

		f.Write([]byte("poyo"))

		files = append(files, f)
	}

	ch := w.OnUpdate()

	name1 := <-ch
	name2 := <-ch

	if name1 == name2 {
		t.Fatal("names should be diffalent, but got same.")
	}

	expected1 := files[0].Name()
	expected2 := files[1].Name()
	if name1 != expected1 && name1 != expected2 {
		t.Fatalf("Name1 should be %s or %s, but got %s", expected1, expected2, name1)
	}
	if name2 != expected1 && name2 != expected2 {
		t.Fatalf("Name2 should be %s or %s, but got %s", expected1, expected2, name2)
	}

	select {
	case n := <-ch:
		t.Errorf("channel should not send data. But got %s", n)
	default:

	}
}
