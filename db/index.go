package db

import(
    "sync"
)

type Index struct{
    paths []string
    indexs map[int][]int
    sync.RWMutex
}


func (idx *Index)Add(){

}
