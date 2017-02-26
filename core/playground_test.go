package core

import (
	"testing"
)

func TestJoin(t *testing.T) {
	p, err := NewPlayground(&Tiki{
		TID:     233,
		Title:   []byte("Hello, Tiki!"),
		Content: []byte("# Intro"),
		Version: 0,
	})

	c := &Client{
		UID: 7,
	}
	_, err = p.Join(c)
	//p.Debug()

	if err != nil {
		t.Error(err.Error())
	}
	p.Left(c)
	//p.Debug()
}
