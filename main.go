package main

import (
	"io/ioutil"
	"os"

	"github.com/google/go-github/github"
)

func main() {
	files := os.Args[1:]

	s := NewServer()

	for _, f := range files {
		md, err := ioutil.ReadFile(f)
		if err != nil {
			s.Add(f, err.Error())
			continue
		}

		html, err := md2html(string(md))
		if err != nil {
			s.Add(f, err.Error())
			continue
		}
		s.Add(f, html)
	}

	select {}
}

func md2html(md string) (string, error) {
	client := github.NewClient(nil)
	html, _, err := client.Markdown(md, nil)
	return html, err
}
