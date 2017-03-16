package core

import (
	"encoding/json"
	"log"
	"sync/core/ot"
)

const (
	OTMsg = iota
	ACKMsg
	ForceSyncMsg
)

// SyncMessage the struct of ot message from fe
type SyncMessage struct {
	Type     uint           `json:"type"`
	UID      uint64         `json:"uid"`
	Sequence uint64         `json:"seq"`
	Version  uint64         `json:"ver"`
	OP       []ot.Operation `json:"op"`
}

// Apply operation to text
func (s *SyncMessage) Apply(content string) string {
	cursor := 0
	for _, op := range s.OP {
		switch op.Type {
		case ot.OPRetain:
			cursor += int(op.Operation.(float64))
		case ot.OPInsert:
			if cursor == len(content) {
				content += op.Operation.(string)
			} else {
				temp := content[:cursor]
				temp += op.Operation.(string)
				temp += content[cursor:]
				cursor += len(op.Operation.(string))
				content = temp
			}
		case ot.OPDelete:
			if len(op.Operation.(string)) == len(content) {
				content = ""
			} else if content[cursor:cursor+len(op.Operation.(string))] == op.Operation.(string) {
				temp := content[:cursor]
				temp += content[cursor+len(op.Operation.(string)):]
				content = temp
			}
		}
	}
	log.Printf("applied %s\n", content)
	return content
}

func (s *SyncMessage) IntentionPreservation(pre *SyncMessage) {
	s.Version += 2
}

// ToJSON marshal the struct to json
func (s *SyncMessage) ToJSON() (msg []byte, err error) {
	msg, err = json.Marshal(s)
	return
}
