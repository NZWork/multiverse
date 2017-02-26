package core

import (
	//"errors"
	"github.com/sergi/go-diff/diffmatchpatch"
	"sync"
)

// Tiki contains everything that TIKI have
type Tiki struct {
	TID     uint64
	Title   []byte
	Content []byte
	Version uint64

	diff *diffmatchpatch.DiffMatchPatch
	lock sync.RWMutex
}

// GetTikiByID will return a tiki instance by ID
func GetTikiByID(id uint64) (*Tiki, error) {
	// Load tiki from storage service
	return &Tiki{
		diff: diffmatchpatch.New(),
	}, nil
}

// GetTikiByToken return a tiki instance by token
func GetTikiByToken(token string) (t *Tiki, err error) {
	//err = errors.New("invalid token")
	t = &Tiki{
		TID: 233,
	}
	return
}
