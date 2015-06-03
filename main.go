package main

import (
	"io/ioutil"
	"os"
)

func main() {
	files := os.Args[1:]

	s := NewServer()

	for _, f := range files {
		md, err := ioutil.ReadFile(f)
		if err != nil {
			s.storage.Add(f, err.Error())
			continue
		}

		html, err := s.storage.md2html(string(md))
		if err != nil {
			s.storage.Add(f, err.Error())
			continue
		}
		s.storage.Add(f, html)
	}

	select {}
}
