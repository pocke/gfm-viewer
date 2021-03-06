package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/naoina/denco"
)

func TestLoadAce(t *testing.T) {
	w := httptest.NewRecorder()
	loadAce(w, "before_auth", nil)

	if w.Code != http.StatusOK {
		t.Errorf("status should be %d, but got %d", http.StatusOK, w.Code)
	}

	w = httptest.NewRecorder()
	loadAce(w, "not_exist_template", nil)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status should be %d, but got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestServeAsset(t *testing.T) {
	s := &Server{}

	fn := func(typ, fname, expectedCtype string) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.serveAsset(w, r, denco.Params{
				denco.Param{Name: "type", Value: typ},
				denco.Param{Name: "fname", Value: fname},
			})
		}))
		defer ts.Close()

		res, err := http.Get(ts.URL)
		if err != nil {
			t.Fatal(err)
		}
		ctype := res.Header.Get("Content-Type")
		if ctype != expectedCtype {
			t.Fatalf("Content-Type should be %s, but got %s", expectedCtype, ctype)
		}
	}

	fn("js", "main.js", "application/javascript")
	fn("css", "bootstrap.min.css", "text/css")
}
