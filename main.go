package main

import (
	"flag"
	"log"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 1124, "TCP port number")
	flag.Parse()

	files := flag.Args()

	s := NewServer(port)

	s.storage.AddFiles(files)

	select {}
}

func Log(format string, args ...interface{}) {
	if DEBUG {
		log.Printf(format, args...)
	}
}
