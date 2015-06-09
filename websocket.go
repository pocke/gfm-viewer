package main

import (
	"net/http"

	"golang.org/x/net/websocket"

	"github.com/naoina/denco"
)

type signal struct{}

type WSManager interface {
	ServeWS(w http.ResponseWriter, r *http.Request, _ denco.Params)
}

func NewWSManager(ch <-chan string) WSManager {
	w := &wsManager{
		onUpdate: ch,
		sessions: make(map[int]chan string),
	}

	go w.watch()

	return w
}

type wsManager struct {
	onUpdate <-chan string
	// TODO: mutex
	sessions map[int]chan string
}

func (wsm *wsManager) ServeWS(w http.ResponseWriter, r *http.Request, _ denco.Params) {
	websocket.Handler(func(ws *websocket.Conn) {
		id := uniqID()
		ch := make(chan string)
		wsm.sessions[id] = ch
		onClose := wsm.onClose(ws)
		for {
			select {
			case path := <-ch:
				ws.Write([]byte(path))
			case <-onClose:
				delete(wsm.sessions, id)
				close(ch)
				return
			}
		}
	}).ServeHTTP(w, r)
}

func (w *wsManager) onClose(ws *websocket.Conn) <-chan signal {
	ch := make(chan signal)
	go func() {
		for {
			_, err := ws.Read(make([]byte, 512))
			if err != nil {
				ch <- signal{}
				break
			}
		}
	}()
	return ch
}

func (w *wsManager) watch() {
	for path := range w.onUpdate {
		for _, ch := range w.sessions {
			go func() { ch <- path }()
		}
	}
}

var uniqID = func() func() int {
	ch := make(chan int)
	go func() {
		i := 0
		for {
			ch <- i
			i++
		}
	}()

	return func() int {
		return <-ch
	}
}()
