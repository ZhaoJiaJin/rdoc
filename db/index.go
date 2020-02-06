package db

import(
    "sync"
    "strings"
)

//Index index struct
type Index struct{
    paths []string
    indexs map[int][]string
    sync.RWMutex
}

//NewIndex create new index
func NewIndex(path string)(*Index){
    return &Index{
        paths:strings.Split(path,","),
        indexs:make(map[int][]string),
    }
}


func (idx *Index)IndexDoc(id string,d *Doc){

}


