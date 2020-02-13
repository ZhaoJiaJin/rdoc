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
	indexs map[int]*IDList
	sync.RWMutex
}

//NewIndex create new index
func NewIndex(path string) *Index {
	return &Index{
		paths:  strings.Split(path, ","),
		indexs: make(map[int]*IDList),
	}
}

//IndexDoc build index for a doc
func (idx *Index) IndexDoc(id string, d *Doc) {
	for _, idxVal := range GetIn(d.data, idx.paths) {
		if idxVal != nil {
			hashKey := utils.StrHash(fmt.Sprint(idxVal))
			idx.Lock()
			if _, ok := idx.indexs[hashKey]; !ok {
				idx.indexs[hashKey] = NewIDList()
			}
			idx.Unlock()
            idx.indexs[hashKey].Add(id)
		}
	}
}

//UnIndex remove id from index
func (idx *Index)UnIndex(id string){
    idx.RLock()
    for _, v := range idx.indexs{
        v.Remove(id)
    }
    idx.RUnlock()
}

//Query index query
func (idx *Index)Query(val int, limit int)([]string,error){
    var ret []string
    idx.RLock();
    idlist,ok := idx.indexs[val];
    idx.RUnlock();
    if !ok{
        return ret,nil
    }
    return idlist.Get(limit),nil
}


//QueryExist index query existence
func (idx *Index)QueryExist(limit int)([]string,error){
    var ret []string
    idx.RLock();
    defer idx.RUnlock();
    //idlist,ok := idx.indexs[val];
    for _,idlist := range idx.indexs{
        tmpres := idlist.Get(limit - len(ret))
        ret = append(ret,tmpres...)
        if len(ret) >= limit{
            return ret,nil
        }
    }
    return ret, nil
}
