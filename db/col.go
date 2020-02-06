package db

import (
	"sync"
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

func (c *Col) UpdateDoc(data []byte) error {
	return nil
}
func (c *Col) MergeDoc(data []byte) error {
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

func (c *Col) CreateIndex(paths string) error {
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
