package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/websocket"

	"github.com/naoina/denco"
)

func TestWSManagerAdd(t *testing.T) {
	ch := make(chan string)
	wsm := NewWSManager(ch)
	path := "poyo"

	wsm.add(path)

	if l := len(wsm.sessions[path]); l != 1 {
		t.Errorf("Expected: 1, but got %d", l)
	}
}

func TestWSManagerWatch(t *testing.T) {
	updateCh := make(chan string)
	wsm := NewWSManager(updateCh)
	path := "poyo"

	watchCh := wsm.add(path)

	updateCh <- path

	<-watchCh

	updateCh <- path // as sleep
	if l := len(wsm.sessions[path]); l != 0 {
		t.Errorf("Expected: 0, but got %d", l)
	}
}

func TestWSManagerServeWS(t *testing.T) {
	ch := make(chan string)
	wsm := NewWSManager(ch)
	path := "poyo"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wsm.ServeWS(w, r, denco.Params{denco.Param{
			Name:  "path",
			Value: path,
		}})
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	ws, err := websocket.Dial(strings.Replace(ts.URL, "http://", "ws://", 1), "", ts.URL)
	if err != nil {
		t.Error(err)
	}

	ch <- path

	buf := make([]byte, 7)
	_, err = ws.Read(buf)
	if err != nil {
		t.Error(err)
	}
	if string(buf) != "Update!" {
		t.Errorf("Expected: Update!, but got %q", buf)
	}
}
