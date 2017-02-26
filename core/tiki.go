package core

import (
	//"errors"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"sync"
)

// Tiki contains everything that TIKI have
type Tiki struct {
	TID     uint64
	Title   []byte
	Content string
	Version uint64

	dmp  *diffmatchpatch.DiffMatchPatch
	lock sync.RWMutex
}

// GetTikiByID will return a tiki instance by ID
func GetTikiByID(id uint64) (*Tiki, error) {
	// Load tiki from storage service
	return &Tiki{
		dmp:     diffmatchpatch.New(),
		Content: "",
		Version: 0,
	}, nil
}

// GetTikiByToken return a tiki instance by token
func GetTikiByToken(token string) (t *Tiki, err error) {
	//err = errors.New("invalid token")
	t = &Tiki{
		TID:     233,
		dmp:     diffmatchpatch.New(),
		Content: "",
		Version: 0,
	}
	return
}

// Merge return merged content
func (t *Tiki) Merge(latest string) {
	diffs := t.dmp.DiffMain(t.Content, latest, true)
	patches := t.dmp.PatchMake(t.Content, diffs)
	t.Content, _ = t.dmp.PatchApply(patches, t.Content)
	t.Version++
	log.Println("merged")
}
