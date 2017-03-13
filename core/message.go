package core

import (
	"encoding/json"
	"sync/core/ot"
)

// SyncMessage the struct of ot message from fe
type SyncMessage struct {
	UID      uint64    `json:"uid"`
	Sequence uint64    `json:"seq"`
	Version  uint64    `json:"ver"`
	OP       Operation `json:"op"`
}

type Operation struct {
	Insert ot.Insert `json:"insert"`
	Delete ot.Delete `json:"delete"`
	Retain ot.Retain `json:"retain"`
}

func (o *Operation) Apply(content string) (newContent string) {
	if o.Retain.Length != 0 {
		newContent = content[:o.Retain.Length]
	}
	if len(o.Insert.Insert) != 0 {

	}
	return
}

// ToJSON marshal the struct to json
func (s *SyncMessage) ToJSON() (msg []byte, err error) {
	msg, err = json.Marshal(s)
	return
}
