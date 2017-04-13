package data

import (
	"log"
	"testing"
)

const Token = ""

func TestRedis(t *testing.T) {
	SetupRedis()
	log.Println("redis init")
	content, err := GetContent("1_123321")
	log.Printf("%s\n", content)
	if err != nil {
		t.Error(err)
	}
	//
	// err = SetContent("1_123321", []byte("fuckme"))
	// if err != nil {
	// 	t.Error(err)
	// }
	//
	// content, err = GetContent("1_123321")
	// log.Printf("%s\n", content)
	// if err != nil {
	// 	t.Error(err)
	// }
}
