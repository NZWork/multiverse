package core

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Tranx for transfer every connect to playground wich it belongs
type Tranx struct {
	playgrounds map[string]*Playground // token => Playground
}

// NewTranx return a Tranx instance
func NewTranx() *Tranx {
	return &Tranx{
		playgrounds: make(map[string]*Playground),
	}
}

// Fuck the world called when new connection created
func (t *Tranx) Fuck(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panic(err.Error())
	}

	p, err := t.getPlayground(r.URL.Query().Get("token"))
	if err != nil { // fail
		log.Printf("connect closed due to %s\n", err.Error())
		return
	}
	c := GetClient(conn, p)
	_, err = p.Join(c)
	if err != nil {
		log.Printf("error when join playground %v\n", err.Error())
		return
	}
	go c.write()
	c.read()
}

func (t *Tranx) getPlayground(token string) (p *Playground, err error) {
	var ok bool
	p, ok = t.playgrounds[token]
	if ok {
		// playground exists
		return
	}
	// playground did not exists or invalid
	p, err = GetPlayground(token)
	if err != nil {
		return
	}
	go p.Run()
	t.playgrounds[token] = p
	return
}
