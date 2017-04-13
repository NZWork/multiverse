package data

import (
	"log"
	"testing"
)

func TestDeserializeToken(t *testing.T) {
	result, err := DeserializeToken("5YoHDAIpComZt5qfz5TEzI6SizfaAjIsKQnx/irp04aWu8SOCfoDjs/3jYgK0wIBh6Sfz5jF9svJ+f3HFQ==",
		"530f6cc1-29480177-2c62ebd3-daa630d5")
	if err != nil {
		t.Error(err)
	}
	log.Printf("%v\n", result)
}

func TestPersistenceTiki(t *testing.T) {
	err := PersistenceTiki(1, "123321")
	if err != nil {
		t.Error(err)
	}
}
