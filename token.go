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

// Token is Github access token.
type Token struct {
	Token string `json:"token"`
}

// Init gets access token from GitHub.
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

// SaveFile saves token to file.
func (t *Token) SaveFile() error {
	return ioutil.WriteFile(t.filePath(), []byte(t.Token), 0644)
}

// LoadFile loads token from file.
func (t *Token) LoadFile() error {
	f, err := ioutil.ReadFile(t.filePath())
	if err != nil {
		return err
	}
	t.Token = string(f)
	return nil
}

// hasToken returns whether the token is saved.
func (t *Token) hasToken() bool {
	err := t.LoadFile()
	if err != nil {
		return false
	}
	return t.Token != ""
}

// filePath is path of saving token.
func (_ *Token) filePath() string {
	fname := "gfm-viewer"
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return path.Join(xdg, fname)
	} else {
		return path.Join(os.Getenv("HOME"), ".cache", fname)
	}
}

// For Storage#md2html
func (t *Token) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "token "+t.Token)

	return http.DefaultTransport.RoundTrip(req)
}
