package db

import (
	"sync"
    "strings"
    "encoding/json"
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
        c.idxlock.RLock()
        for _,idxe := range c.index{
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
        delete(c.docs,id)
        c.doclock.Unlock()
        c.idxlock.RLock()
        for _,idxe := range c.index{
            idxe.UnIndex(id)
        }
        c.idxlock.RUnlock()
    }
	return nil
}

//QueryDocID query document using index
func (c *Col) QueryDocID(data []byte) (res map[string]struct{}, err error) {
    var qJSON interface{}
	if err = json.Unmarshal(data, &qJSON); err != nil {
		//http.Error(w, fmt.Sprintf("'%v' is not valid JSON.", q), 400)
		return
	}
    err = EvalQuery(qJSON,c,&res)
    return
}

//ReadDoc read document by id
func (c *Col) ReadDoc(id string) (res *Doc) {
    c.doclock.RLock()
    res = c.docs[id]
    c.doclock.RUnlock()
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

//RmIndex remove index
func (c *Col) RmIndex(paths string) error {
    c.idxlock.Lock()
    delete(c.index,paths)
    c.idxlock.Unlock()
	return nil
}

//GetAllIndex get all index
func (c *Col) GetAllIndex() (ret []string, err error) {
    c.idxlock.RLock()
    for k := range c.index{
        ret = append(ret,k)
    }
    c.idxlock.RUnlock()
	return
}

//ForEachDoc iterate all documents
func (c *Col) ForEachDoc(fun func(id string, doc *Doc) (bool)) {
    c.doclock.RLock()
    defer c.doclock.RUnlock()
    for k,v := range c.docs{
        if !fun(k,v){
            return
        }
    }
}

// Query query using index
func (c *Col)Query(scanpath string, lookupValueHash int, limit int)([]string,error){
    var ret []string
    c.idxlock.RLock()
    pathidx,ok := c.index[scanpath]
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
    pathidx,ok := c.index[scanpath]
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

    if _, ok := c.index[path]; !ok{
        return false
    }
    return true
}
