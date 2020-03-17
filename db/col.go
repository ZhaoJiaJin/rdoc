package db

import (
    "strings"
    "encoding/json"
)

//Col collection
type Col struct {
	Docs    map[string]*Doc
	Index   map[string]*Index //what if there are duplicated index entry
}

//NewCol Create New Collection
func NewCol() *Col {
	return &Col{
		Docs:  make(map[string]*Doc),
		Index: make(map[string]*Index),
	}
}

//AddDoc add doc into collection
func (c *Col) AddDoc(id string,data []byte) (string, error) {
	doc, err := NewDoc(data)
	if err != nil {
		return id, err
	}
	c.Docs[id] = doc

	for _, idx := range c.Index {
		idx.IndexDoc(id, doc)
	}
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
        c.Docs[id] = newdoc
        for _,idxe := range c.Index{
            idxe.UnIndex(id)
            idxe.IndexDoc(id,newdoc)
        }
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
        if _,ok := c.Docs[id]; !ok{
            continue
        }
        c.Docs[id].Merge(newdoc)
        for _,idxe := range c.Index{
            idxe.UnIndex(id)
            idxe.IndexDoc(id,c.Docs[id])
        }
    }
	return nil
}

//DeleteDoc delete documents
func (c *Col) DeleteDoc(ids string) error {
    for _,id := range strings.Split(ids,","){
        delete(c.Docs,id)
        for _,idxe := range c.Index{
            idxe.UnIndex(id)
        }
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
    res = c.Docs[id]
	return
}

//CreateIndex create index for a collection
func (c *Col) CreateIndex(paths string) error {
    if _,ok := c.Index[paths]; ok{
        return ErrIDXExist
    }

    nidx := NewIndex(paths)
    c.Index[paths] = nidx

    for id,doc := range c.Docs{
        nidx.IndexDoc(id,doc)
    }
	return nil
}

//RmIndex remove index
func (c *Col) RmIndex(paths string) error {
    delete(c.Index,paths)
	return nil
}

//GetAllIndex get all index
func (c *Col) GetAllIndex() (ret []string, err error) {
    for k := range c.Index{
        ret = append(ret,k)
    }
	return
}

//ForEachDoc iterate all documents
func (c *Col) ForEachDoc(fun func(id string, doc *Doc) (bool)) {
    for k,v := range c.Docs{
        if !fun(k,v){
            return
        }
    }
}

// Query query using index
func (c *Col)Query(scanpath string, lookupValueHash int, limit int)([]string,error){
    var ret []string
    pathidx,ok := c.Index[scanpath]
    if !ok{
        return ret, ErrNotIDX
    }
    return pathidx.Query(lookupValueHash,limit)
}

// QueryExist query existence using index
func (c *Col)QueryExist(scanpath string, limit int)([]string,error){
    var ret []string
    pathidx,ok := c.Index[scanpath]
    if !ok{
        return ret, ErrNotIDX
    }
    return pathidx.QueryExist(limit)
}

func (c *Col)IsIndexed(path string)bool{

    if _, ok := c.Index[path]; !ok{
        return false
    }
    return true
}
