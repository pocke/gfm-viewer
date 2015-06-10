package main

import (
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type signal struct{}

type WSManager interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func NewWSManager(ch <-chan string) WSManager {
	w := &wsManager{
		onUpdate: ch,
		sessions: make(map[int]chan string),
		mu:       &sync.RWMutex{},
	}

	go w.watch()

	return w
}

type wsManager struct {
	onUpdate <-chan string

	sessions map[int]chan string
	mu       *sync.RWMutex
}

func (wsm *wsManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	websocket.Handler(func(ws *websocket.Conn) {
		id := uniqID()
		ch := make(chan string)
		wsm.withLock(func() { wsm.sessions[id] = ch })
		onClose := wsm.onClose(ws)
		for {
			select {
			case path := <-ch:
				ws.Write([]byte(path))
			case <-onClose:
				wsm.withLock(func() { delete(wsm.sessions, id) })
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
		w.mu.RLock()
		for _, ch := range w.sessions {
			go func() { ch <- path }()
		}
		w.mu.RUnlock()
	}
}

func (w *wsManager) withLock(f func()) {
	w.mu.Lock()
	defer w.mu.Unlock()
	f()
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
