package core

import (
	"log"
	"multiverse/data"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
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
	token := r.URL.Query().Get("token")
	pubkey := r.URL.Query().Get("pubkey")

	log.Printf("token %s, pubkey %s\n", token, pubkey)

	result, err := data.DeserializeToken(token, pubkey)
	if err != nil {
		log.Printf("failed to deserialize token and pubkey due to %s\n", err.Error())
		conn.Close()
		return
	}

	p, err := t.getPlayground(result.ProjectID + "_" + result.Token)
	log.Printf("Project ID %v token %v\n", result.ProjectID, result.Token)
	if err != nil { // fail
		log.Printf("connect closed due to %s\n", err.Error())
		conn.Close()
		return
	}
	c := GetClient(conn, p)
	c.UID = uint64(result.UID)

	_, err = p.Join(c)
	if err != nil {
		log.Printf("error when join playground %v\n", err.Error())
		conn.Close()
		return
	}
	go c.write()
	// Sync tiki content
	p.InitClient(c) // init client
	// p.ChaseTiki(c)  // chase the progress of tiki
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
	p, err = GetPlayground(token, t)
	if err != nil {
		return
	}
	go p.Run()
	t.playgrounds[token] = p
	return
}

func (t *Tranx) closePlayground(token string) {
	temp := strings.Split(token, "_")
	if len(temp) == 2 {
		err := data.PersistenTiki(temp[0], temp[1])
		if err != nil {
			log.Printf("error when persisten tiki: %s", err.Error())
		}
	}
	delete(t.playgrounds, token)
}
