package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestSaveFileLoadFile(t *testing.T) {
	reset := helperSetEnv()
	defer reset()

	str := "hogepoyoyo"
	token := &Token{Token: str}

	err := token.SaveFile()
	if err != nil {
		t.Error(err)
	}

	token = &Token{}
	err = token.LoadFile()
	if err != nil {
		t.Error(err)
	}

	if token.Token != str {
		t.Errorf("Expected: %s, but got %s", str, token.Token)
	}
}

func TestHasToken(t *testing.T) {
	reset := helperSetEnv()
	defer reset()

	token := &Token{}
	if token.hasToken() {
		t.Errorf("Expected: false, but got true")
	}

	token.Token = "hoge"
	token.SaveFile()

	token = &Token{}
	if !token.hasToken() {
		t.Errorf("Expected: true, but got false")
	}
}

func helperSetEnv() func() {
	key := "XDG_CACHE_HOME"
	pre := os.Getenv(key)

	name, err := ioutil.TempDir("", "gfm-viewer-test")
	if err != nil {
		panic(err)
	}

	os.Setenv(key, name)
	return func() {
		os.RemoveAll(name)
		os.Setenv(key, pre)
	}
}
