package db

import(
    "sync"
)

type Col struct{
    docs map[int]map[string]interface{}
    index map[string]*Index //what if there are duplicated index entry
    sync.RWMutex
}


func (c *Col)AddDoc(data []byte)(error){
    d,err := NewDoc(data)
    if err != nil{
        return err
    }
}



func (c *Col)CreateIndex(paths []string){

}


// rebuildIndex should be called everytime a new index is added
func (c *Col)rebuildIndex(){

}
