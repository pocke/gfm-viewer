package main

import (
	"net/http"

	"github.com/naoina/denco"
	"golang.org/x/net/websocket"
)

type WSManager struct {
	sessions      map[string][]chan signal
	receiveUpdate <-chan string
}

type signal struct{}

func NewWSManager(ch <-chan string) *WSManager {
	w := &WSManager{
		sessions:      make(map[string][]chan signal),
		receiveUpdate: ch,
	}
	go w.watch()

	return w
}

func (wsm *WSManager) ServeWS(w http.ResponseWriter, r *http.Request, p denco.Params) {
	websocket.Handler(func(ws *websocket.Conn) {
		ch := wsm.add(p.Get("path"))
		<-ch
		ws.Write([]byte("Update!"))
	}).ServeHTTP(w, r)
}

func (w *WSManager) add(path string) <-chan signal {
	s, ok := w.sessions[path]
	if !ok {
		s = make([]chan signal, 0, 1)
	}

	ch := make(chan signal)
	w.sessions[path] = append(s, ch)
	return ch
}

func (w *WSManager) watch() {
	for {
		path := <-w.receiveUpdate
		Log("Update %s", path)
		s, ok := w.sessions[path]
		if !ok {
			continue
		}
		for _, v := range s {
			v <- signal{}
		}
		delete(w.sessions, path)
	}
}
