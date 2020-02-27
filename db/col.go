package db

import (
	"sync"
    "strings"
    "encoding/json"
)

//Col collection
type Col struct {
	Docs    map[string]*Doc
	Index   map[string]*Index //what if there are duplicated index entry
	doclock sync.RWMutex
	idxlock sync.RWMutex
}

//NewCol Create New Collection
func NewCol() *Col {
	return &Col{
		Docs:  make(map[string]*Doc),
		Index: make(map[string]*Index),
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
	c.Docs[id] = doc
	c.doclock.Unlock()

	c.idxlock.RLock()
	for _, idx := range c.Index {
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
        c.Docs[id] = newdoc
        c.doclock.Unlock()
        c.idxlock.RLock()
        for _,idxe := range c.Index{
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
        if _,ok := c.Docs[id]; !ok{
            c.doclock.Unlock()
            continue
        }
        c.Docs[id].Merge(newdoc)
        c.doclock.Unlock()
        c.idxlock.RLock()
        for _,idxe := range c.Index{
            idxe.UnIndex(id)
            idxe.IndexDoc(id,newdoc)
        }
        c.idxlock.RUnlock()
    }
	return nil
}

//DeleteDoc delete documents
func (c *Col) DeleteDoc(ids string) error {
    for _,id := range strings.Split(ids,","){
        c.doclock.Lock()
        delete(c.Docs,id)
        c.doclock.Unlock()
        c.idxlock.RLock()
        for _,idxe := range c.Index{
            idxe.UnIndex(id)
        }
        c.idxlock.RUnlock()
    }
	return nil
}

//QueryDocID query document using index
func (c *Col) QueryDocID(data []byte) (res map[string]struct{}, err error) {
    res = make(map[string]struct{})
    var qJSON interface{}
	if err = json.Unmarshal(data, &qJSON); err != nil {
		return
	}
    err = EvalQuery(qJSON,c,&res)
    return
}

//ReadDoc read document by id
func (c *Col) ReadDoc(id string) (res *Doc) {
    c.doclock.RLock()
    res = c.Docs[id]
    c.doclock.RUnlock()
	return
}

//CreateIndex create index for a collection
func (c *Col) CreateIndex(paths string) error {
    c.idxlock.RLock()
    if _,ok := c.Index[paths]; ok{
        c.idxlock.RUnlock()
        return ErrIDXExist
    }
    c.idxlock.RUnlock()

    c.idxlock.Lock()
    nidx := NewIndex(paths)
    c.Index[paths] = nidx
    c.idxlock.Unlock()

    c.doclock.RLock()
    for id,doc := range c.Docs{
        nidx.IndexDoc(id,doc)
    }
    c.doclock.RUnlock()
	return nil
}

//RmIndex remove index
func (c *Col) RmIndex(paths string) error {
    c.idxlock.Lock()
    delete(c.Index,paths)
    c.idxlock.Unlock()
	return nil
}

//GetAllIndex get all index
func (c *Col) GetAllIndex() (ret []string, err error) {
    c.idxlock.RLock()
    for k := range c.Index{
        ret = append(ret,k)
    }
    c.idxlock.RUnlock()
	return
}

//ForEachDoc iterate all documents
func (c *Col) ForEachDoc(fun func(id string, doc *Doc) (bool)) {
    c.doclock.RLock()
    defer c.doclock.RUnlock()
    for k,v := range c.Docs{
        if !fun(k,v){
            return
        }
    }
}

// Query query using index
func (c *Col)Query(scanpath string, lookupValueHash int, limit int)([]string,error){
    var ret []string
    c.idxlock.RLock()
    pathidx,ok := c.Index[scanpath]
    if !ok{
        c.idxlock.RUnlock()
        return ret, ErrNotIDX
    }
    c.idxlock.RUnlock()
    return pathidx.Query(lookupValueHash,limit)
}

// QueryExist query existence using index
func (c *Col)QueryExist(scanpath string, limit int)([]string,error){
    var ret []string
    c.idxlock.RLock()
    pathidx,ok := c.Index[scanpath]
    if !ok{
        c.idxlock.RUnlock()
        return ret, ErrNotIDX
    }
    c.idxlock.RUnlock()
    return pathidx.QueryExist(limit)
}

func (c *Col)IsIndexed(path string)bool{
    c.idxlock.RLock()
    defer c.idxlock.RUnlock()

    if _, ok := c.Index[path]; !ok{
        return false
    }
    return true
}
