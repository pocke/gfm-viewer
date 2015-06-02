package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-github/github"
)

func main() {
	files := os.Args[1:]

	for _, f := range files {
		md, err := ioutil.ReadFile(f)
		if err != nil {
			continue
		}

		html, err := md2html(string(md))
		if err != nil {
			continue
		}
		fmt.Print(html)
	}
}

func md2html(md string) (string, error) {
	client := github.NewClient(nil)
	html, _, err := client.Markdown(md, nil)
	return html, err
}
