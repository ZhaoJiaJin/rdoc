package db

import (
	"sync"
    "strings"
)

//Col collection
type Col struct {
	docs    map[string]*Doc
	index   map[string]*Index //what if there are duplicated index entry
	doclock sync.RWMutex
	idxlock sync.RWMutex
}

//NewCol Create New Collection
func NewCol() *Col {
	return &Col{
		docs:  make(map[string]*Doc),
		index: make(map[string]*Index),
	}
}

//AddDoc add doc into collection
func (c *Col) AddDoc(data []byte) (string, error) {
	id := RandID()
	doc, err := NewDoc(data)
	if err != nil {
		return id, err
	}
	c.doclock.Lock()
	c.docs[id] = doc
	c.doclock.Unlock()

	c.idxlock.RLock()
	for _, idx := range c.index {
		idx.IndexDoc(id, doc)
	}
	c.idxlock.RUnlock()
	//generate id
	return id, nil
}

//UpdateDoc update document
func (c *Col) UpdateDoc(ids string, data []byte) error {
    for _,id := range strings.Split(ids,","){
        newdoc,err := NewDoc(data)
        if err != nil{
            return err
        }
        c.doclock.Lock()
        c.docs[id] = newdoc
        c.doclock.Unlock()
        //TODO
        c.idxlock.RLock()
        for _,idxe := range c.index{
            idxe.UnIndex(id)
            idxe.IndexDoc(id,newdoc)
        }
        c.idxlock.RUnlock()
    }
	return nil
}

//MergeDoc merge given document to existing document
func (c *Col) MergeDoc(ids string, data []byte) error {
    newdoc,err := NewDoc(data)
    if err != nil{
        return err
    }
    for _,id := range strings.Split(ids,","){
        c.doclock.Lock()
        c.docs[id].Merge(newdoc)
        c.doclock.Unlock()
        //TODO
        c.idxlock.RLock()
        for _,idxe := range c.index{
            idxe.UnIndex(id)
            idxe.IndexDoc(id,newdoc)
        }
        c.idxlock.RUnlock()
    }
	return nil
}
func (c *Col) DeleteDoc(ids string) error {
	return nil
}
func (c *Col) QueryDocID(data []byte) (res []string, err error) {
	return
}
func (c *Col) ReadDoc(id string) (res *Doc) {
	return
}

//CreateIndex create index for a collection
func (c *Col) CreateIndex(paths string) error {
    c.idxlock.RLock()
    if _,ok := c.index[paths]; ok{
        c.idxlock.RUnlock()
        return ErrIDXExist
    }
    c.idxlock.RUnlock()

    c.idxlock.Lock()
    nidx := NewIndex(paths)
    c.index[paths] = nidx
    c.idxlock.Unlock()

    c.doclock.RLock()
    for id,doc := range c.docs{
        nidx.IndexDoc(id,doc)
    }
    c.doclock.RUnlock()
	return nil
}

func (c *Col) RmIndex(paths string) error {
	return nil
}

func (c *Col) GetAllIndex() (ret []string, err error) {

	return
}

// rebuildIndex should be called everytime a new index is added
func (c *Col) rebuildIndex() {

}
