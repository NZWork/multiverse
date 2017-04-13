package core

import (
	"log"
	"multiverse/data"
	"strings"
	"sync"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// Tiki contains everything that TIKI have
type Tiki struct {
	ProjectID  string
	Token      string
	CacheToken string
	Content    string
	Version    uint64

	dmp  *diffmatchpatch.DiffMatchPatch
	lock sync.RWMutex
}

// GetTiki return a tiki instance by token
func GetTiki(token string) (t *Tiki, err error) {
	// split token
	temp := strings.Split(token, "_")
	content, err := data.GetContent(token)
	if err != nil {
		return nil, err
	}

	t = &Tiki{
		ProjectID:  temp[0],
		Token:      temp[1],
		dmp:        diffmatchpatch.New(),
		Content:    string(content),
		CacheToken: token,
		Version:    0,
	}
	return
}

// UpdateCache save content
func (t *Tiki) UpdateCache() (err error) {
	err = data.SetContent(t.CacheToken, []byte(t.Content))
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
