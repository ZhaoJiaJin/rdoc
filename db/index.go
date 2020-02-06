package db

import(
    "sync"
    "strings"
)

//Index index struct
type Index struct{
    paths []string
    indexs map[int][]*ID
    sync.RWMutex
}

//NewIndex create new index
func NewIndex(path string)(*Index){
    return &Index{
        path:strings.Split(path,","),
        indexs:make(map[int][]*ID),
    }
}


func (idx *Index)IndexDoc(id *ID,d *Doc){

}


