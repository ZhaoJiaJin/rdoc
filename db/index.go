package db

import (
	"fmt"
	"strings"
	"sync"
    "rdoc/utils"
)

//Index index struct
type Index struct {
	paths  []string
	indexs map[int][]string
	sync.RWMutex
}

//NewIndex create new index
func NewIndex(path string) *Index {
	return &Index{
		paths:  strings.Split(path, ","),
		indexs: make(map[int][]string),
	}
}

//IndexDoc build index for a doc
func (idx *Index) IndexDoc(id string, d *Doc) {
	for _, idxVal := range GetIn(d.data, idx.paths) {
		if idxVal != nil {
			hashKey := utils.StrHash(fmt.Sprint(idxVal))
			idx.Lock()
			if _, ok := idx.indexs[hashKey]; !ok {
				idx.indexs[hashKey] = make([]string, 0)
			}
            idx.indexs[hashKey] = append(idx.indexs[hashKey],id)
			idx.Unlock()
		}
	}
}


//func (idx *Index)Query()
