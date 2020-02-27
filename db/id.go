package db

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
    "container/list"
    "sync"
    "encoding/json"
)

var (
	node *snowflake.Node
)

//IDList idlist structure
type IDList struct{
    list *list.List
    sync.RWMutex
}

//NewIDList create new idlist
func NewIDList()*IDList{
    return &IDList{
        list:list.New(),
    }
}

func (l *IDList)UnmarshalJSON(b []byte) error {
    l.Lock()
    defer l.Unlock()
    res := make([]string,0)
    if err := json.Unmarshal(b, &res); err != nil{
        return err
    }
    l.list = list.New()
    for _,v := range res{
        l.list.PushBack(v)
    }
    return nil
}

func (l *IDList)MarshalJSON() ([]byte, error) {
    l.Lock()
    defer l.Unlock()
    res := make([]string, l.list.Len())
    idx := 0
    for e := l.list.Front(); e != nil; e = e.Next() {
        res[idx] = e.Value.(string)
        idx ++
	}
    return json.Marshal(res)
}

//Add add element to list
func (l *IDList)Add(id string){
    l.Lock()
    l.list.PushBack(id)
    l.Unlock()
}

//Remove remove element from list
func (l *IDList)Remove(id string){
    l.RLock()
    ret := make([]*list.Element,0)
    for e := l.list.Front(); e != nil; e = e.Next() {
        idstr := e.Value.(string)
        if (idstr == id){
            ret = append(ret,e)
        }
	}
    l.RUnlock()
    l.Lock()
    for _,e := range ret{
        l.list.Remove(e)
    }
    l.Unlock()
}

//Get get all id
func (l *IDList)Get(limit int)(ret []string){
    cnt := 0
    l.RLock()
    defer l.RUnlock()
    for e := l.list.Front(); e != nil; e = e.Next() {
        ret = append(ret,e.Value.(string))
        cnt ++
        if cnt == limit{
            return
        }
	}
    return
}

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

//RandID generate randomid
func RandID() string {
	return fmt.Sprintf("%X", node.Generate().Int64())
}

