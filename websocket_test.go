package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/websocket"

	"github.com/naoina/denco"
)

var _ WSManager = &wsManager{}

func TestWSManager(t *testing.T) {
	ch := make(chan string)
	wsm := NewWSManager(ch)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wsm.ServeWS(w, r, denco.Params{})
	}))
	defer ts.Close()

	ws, err := websocket.Dial(strings.Replace(ts.URL, "http://", "ws://", 1), "", ts.URL)
	if err != nil {
		t.Error(err)
	}
	defer ws.Close()

	ev := make(chan string)
	go func() {
		b := make([]byte, 512)
		n, err := ws.Read(b)
		if err != nil {
			t.Error(err)
		}
		ev <- string(b[0:n])
	}()

	path := "foobarhoge"
	ch <- path
	got := <-ev

	if path != got {
		t.Errorf("Expected %s, but got %s", path, got)
	}
}
