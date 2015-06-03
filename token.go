package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

type Token struct {
	T string `json:"token"`
}

func (t *Token) Init(user, pass string) error {
	req, err := http.NewRequest(
		"POST",
		"https://api.github.com/authorizations",
		strings.NewReader(`{"note":"gfm-viewer"}`),
	)
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, pass)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf(string(body))
	}

	json.NewDecoder(res.Body).Decode(t)
	return t.SaveFile()
}

func (t *Token) SaveFile() error {
	return ioutil.WriteFile(t.filePath(), []byte(t.T), 0644)
}

func (_ *Token) filePath() string {
	fname := "gfm-viewer"
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return path.Join(xdg, fname)
	} else {
		return path.Join(os.Getenv("HOME"), ".cache", fname)
	}
}

func (t *Token) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "token "+t.T)

	return http.DefaultTransport.RoundTrip(req)
}
