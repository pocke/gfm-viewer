package main

import "os"

func main() {
	files := os.Args[1:]

	s := NewServer()

	s.storage.AddFiles(files)

	select {}
}
