package db

import (
	"fmt"
	"strings"
	"sync"
    "rdoc/utils"
)

//Index index struct
type Index struct {
	Paths  []string
	Indexs map[int]*IDList
	sync.RWMutex
}

//NewIndex create new index
func NewIndex(path string) *Index {
	return &Index{
		Paths:  strings.Split(path, INDEX_PATH_SEP),
		Indexs: make(map[int]*IDList),
	}
}

//IndexDoc build index for a doc
func (idx *Index) IndexDoc(id string, d *Doc) {
	for _, idxVal := range GetIn(d.Data, idx.Paths) {
		if idxVal != nil {
			//hashKey := utils.StrHash(fmt.Sprint(idxVal))
			hashKey := utils.StrHash(fmt.Sprint(idxVal))
			idx.Lock()
			if _, ok := idx.Indexs[hashKey]; !ok {
				idx.Indexs[hashKey] = NewIDList()
			}
			idx.Unlock()
            idx.Indexs[hashKey].Add(id)
		}
	}
}

//UnIndex remove id from index
func (idx *Index)UnIndex(id string){
    idx.RLock()
    for _, v := range idx.Indexs{
        v.Remove(id)
    }
    idx.RUnlock()
}

//Query index query
func (idx *Index)Query(val int, limit int)([]string,error){
    var ret []string
    idx.RLock();
    idlist,ok := idx.Indexs[val];
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
    for _,idlist := range idx.Indexs{
        tmpres := idlist.Get(limit - len(ret))
        ret = append(ret,tmpres...)
        if len(ret) >= limit{
            return ret,nil
        }
    }
    return ret, nil
}
