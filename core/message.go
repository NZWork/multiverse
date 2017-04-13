package core

import (
	"encoding/json"
	"multiverse/core/ot"
)

const (
	OTMsg = iota
	ACKMsg
	ForceSyncMsg
	InitMsg
)

// SyncMessage the struct of ot message from fe
type SyncMessage struct {
	Type     uint         `json:"type"`
	UID      uint64       `json:"uid"`
	Sequence uint64       `json:"seq"`
	Version  uint64       `json:"ver"`
	Change   ot.Changeset `json:"ops"`
}

func Clone(ops []ot.Operation) []ot.Operation {
	return ops
}

// ToJSON marshal the struct to json
func (s *SyncMessage) ToJSON() (msg []byte, err error) {
	msg, err = json.Marshal(s)
	return
}

func min(l1, l2 uint) uint {
	if l1 < l2 {
		return l1
	}
	return l2
}
