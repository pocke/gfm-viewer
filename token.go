package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Token struct {
	Token string `json:"token"`
}

func NewToken(user, pass string) (*Token, error) {
	req, err := http.NewRequest(
		"POST",
		"https://api.github.com/authorizations",
		strings.NewReader(`{"note":"gfm-viewer"}`),
	)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, pass)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf(string(body))
	}

	t := &Token{}
	json.NewDecoder(res.Body).Decode(t)

	return t, nil
}
