package core

import (
	"fmt"
	"sync/core/ot"
)

// DebugSeperator for output
const DebugSeperator = "\n======= %v =======\n"

// Playground for pices of TIKI
type Playground struct {
	Tiki    *Tiki
	Clients map[*Client]bool

	broadcast chan []byte
}

// GetPlayground by token
func GetPlayground(t string) (*Playground, error) {
	// doing things
	tiki, err := GetTikiByToken(t)
	if err != nil {
		return nil, err
	}
	return NewPlayground(tiki)
}

// NewPlayground return a *Playground instance and error
func NewPlayground(t *Tiki) (*Playground, error) {
	fmt.Println("New playground")
	return &Playground{
		Tiki:    t,
		Clients: make(map[*Client]bool),

		broadcast: make(chan []byte),
	}, nil
}

// Join for new client join into playground
func (p *Playground) Join(c *Client) (ok bool, err error) {
	p.Clients[c] = true
	fmt.Printf("UID: %d join\n", c.UID)
	//p.LoadTiki(c)
	return true, nil
}

// Left called when client close the connection
func (p *Playground) Left(c *Client) {
	fmt.Printf("UID: %d left\n", c.UID)
	delete(p.Clients, c)
}

// LoadTiki when the client join
// To sync the lag tiki
func (p *Playground) LoadTiki(c *Client) {
	msg := &SyncMessage{}
	msg.UID = c.UID
	msg.Version = p.Tiki.Version
	msg.Sequence = 0
	msg.OP = Operation{
		Insert: ot.Insert{
			Insert: p.Tiki.Content,
		},
	}
	m, err := msg.ToJSON()
	if err != nil {
		panic(err)
	}
	c.send <- m
}

// Run goroutine handling A playground connections
func (p *Playground) Run() {
	for {
		select {
		case msg := <-p.broadcast:
			fmt.Printf("broadcast: %s\n", msg)
			for c := range p.Clients {
				// 消息发送
				c.send <- msg
			}
		}
	}
}

// Debug for debugging Playground only with ungly performance
func (p *Playground) Debug() {
	// Display TIKI
	fmt.Print(fmt.Sprintf(DebugSeperator, "Tiki"))
	fmt.Printf("TID: %d\n", p.Tiki.TID)
	fmt.Printf("Title: %s\n", p.Tiki.Title)
	fmt.Printf("Content: %s\n", p.Tiki.Content)
	fmt.Printf("Version: %d\n", p.Tiki.Version)
	// Display Clients
	fmt.Print(fmt.Sprintf(DebugSeperator, "Clients"))
	for c := range p.Clients {
		fmt.Printf("ID: %d\n", c.UID)
	}
}
