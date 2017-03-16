package core

import (
	"log"
	"sync/core/ot"
	"testing"
)

func TestIntentionPreservation(t *testing.T) {
	pre := &SyncMessage{
		Type:    OTMsg,
		Version: 3,
		OP: []ot.Operation{
			ot.Delete("1"),
			ot.Retain(1),
			ot.Insert("2"),
			ot.Retain(1),
			ot.Delete("4"),
		},
	}
	now := &SyncMessage{
		Type:    OTMsg,
		Version: 2,
		OP: []ot.Operation{
			ot.Retain(2),
			ot.Delete("3"),
			ot.Retain(1),
			ot.Delete("5"),
			ot.Retain(4),
			ot.Insert("0"),
		},
	}

	now.IntentionPreservation(pre)
	json, _ := now.ToJSON()
	log.Printf("%s", json)
}
