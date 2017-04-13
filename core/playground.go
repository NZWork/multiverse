package core

import (
	"encoding/json"
	"fmt"
	"log"
	"multiverse/core/ot"
)

// DebugSeperator for output
const DebugSeperator = "\n======= %v =======\n"

// Playground for pices of TIKI
type Playground struct {
	Tiki    *Tiki
	Clients map[uint64]*Client

	opHistory map[uint64]*SyncMessage
	tranx     *Tranx
	token     string
	broadcast chan []byte
}

// GetPlayground by token
func GetPlayground(t string, tr *Tranx) (*Playground, error) {
	// doing things
	tiki, err := GetTiki(t)
	if err != nil {
		return nil, err
	}
	return NewPlayground(tiki, t, tr)
}

// NewPlayground return a *Playground instance and error
func NewPlayground(t *Tiki, token string, tranx *Tranx) (*Playground, error) {
	fmt.Println("New playground")
	return &Playground{
		Tiki:    t,
		Clients: make(map[uint64]*Client),

		opHistory: make(map[uint64]*SyncMessage),
		tranx:     tranx,
		token:     token,
		broadcast: make(chan []byte),
	}, nil
}

// Join for new client join into playground
func (p *Playground) Join(c *Client) (ok bool, err error) {
	p.Clients[c.UID] = c
	fmt.Printf("UID: %d join\n", c.UID)
	return true, nil
}

// Left called when client close the connection
func (p *Playground) Left(c *Client) {
	fmt.Printf("UID: %d left\n", c.UID)
	delete(p.Clients, c.UID)
	if len(p.Clients) == 0 {
		p.tranx.closePlayground(p.token)
	}
}

// ChaseTiki when the client join
// To sync the lag tiki
func (p *Playground) ChaseTiki(c *Client) {
	if len(p.Tiki.Content) == 0 {
		return
	}

	length := uint(len(p.Tiki.Content))
	msg := &SyncMessage{}
	msg.Type = ForceSyncMsg
	msg.UID = c.UID
	msg.Version = p.Tiki.Version
	msg.Sequence = 0
	msg.Change = ot.Changeset{
		OP: []ot.Operation{
			ot.Insert(length),
		},
		Adden:        p.Tiki.Content,
		InputLength:  0,
		OutputLength: length,
	}
	m, err := msg.ToJSON()
	if err != nil {
		panic(err)
	}
	c.send <- m
}

// InitClient 初始化客户端
func (p *Playground) InitClient(c *Client) {
	msg := &SyncMessage{
		Type: InitMsg,
		UID:  c.UID,
	}
	m, err := msg.ToJSON()
	if err != nil {
		panic(err)
	}
	c.send <- m
}

// Run goroutine handling A playground connections
func (p *Playground) Run() {
	var err error
	for {
		select {
		case msg := <-p.broadcast:
			// content apply changes
			m := &SyncMessage{}
			json.Unmarshal(msg, m)
			if m.Version == p.Tiki.Version || m.Version == p.Tiki.Version-1 {
				shouldChase := false
				if m.Version == p.Tiki.Version-1 {
					// intention-preservation
					log.Println("intention preservation")
					m.Change.IntentionPreservation(&p.opHistory[m.Version].Change)
					m.Version += 2
					shouldChase = true
				}

				// save to history
				p.opHistory[m.Version] = m
				p.Tiki.Content, err = m.Change.Apply(p.Tiki.Content) // apply to server
				err = p.Tiki.UpdateCache()
				if err != nil {
					log.Printf("failed to save content to redis %s", err.Error())
				} else {
					log.Println("save to redis")
				}
				log.Printf("v%d content %s\n", p.Tiki.Version, p.Tiki.Content)

				if err != nil {
					panic(err)
				}
				m.Version++
				p.Tiki.Version = m.Version

				msg, _ = m.ToJSON()
				log.Printf("broadcast: %s\n", msg)
				for uid, c := range p.Clients {
					if uid == m.UID { // ack
						if shouldChase {
							p.ChaseTiki(c)
							// log.Println("should chase")
						} else {
							temp := &SyncMessage{
								Type:     ACKMsg,
								Version:  m.Version,
								Sequence: m.Sequence,
							}
							ack, _ := temp.ToJSON()
							c.send <- ack
						}
						continue
					}
					// 消息发送
					c.send <- msg
				}
			} else { // this client lag the server version
				// chase
				log.Printf("v%d force sync %d\n", p.Tiki.Version, m.UID)
				p.ChaseTiki(p.Clients[m.UID])
			}
		}
	}
}

// Debug for debugging Playground only with ungly performance
func (p *Playground) Debug() {
	// Display TIKI
	fmt.Print(fmt.Sprintf(DebugSeperator, "Tiki"))
	fmt.Printf("ProjectID: %s\n", p.Tiki.ProjectID)
	fmt.Printf("Token: %s\n", p.Tiki.Token)
	fmt.Printf("Content: %s\n", p.Tiki.Content)
	fmt.Printf("Version: %d\n", p.Tiki.Version)
	// Display Clients
	fmt.Print(fmt.Sprintf(DebugSeperator, "Clients"))
	for _, c := range p.Clients {
		fmt.Printf("ID: %d\n", c.UID)
	}
}
