package main

import (
	"log"
	"os"

	"github.com/pocke/gfm-viewer/env"
)

func main() {
	files := os.Args[1:]

	s := NewServer()

	s.storage.AddFiles(files)

	select {}
}

func Log(format string, args ...interface{}) {
	if env.DEBUG {
		log.Printf(format, args...)
	}
}
