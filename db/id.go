package db

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
    "container/list"
)

var (
	node *snowflake.Node
)

//IDList idlist structure
type IDList struct{
    list *list.List
}

//NewIDList create new idlist
func NewIDList()*IDList{
    return &IDList{
        list:list.New(),
    }
}

//Add add element to list
func (l *IDList)Add(id string){
    l.list.PushBack(id)
}

//Remove remove element from list
func (l *IDList)Remove(id string){
    ret := make([]*list.Element,0)
    for e := l.list.Front(); e != nil; e = e.Next() {
        idstr := e.Value.(string)
        if (idstr == id){
            ret = append(ret,e)
        }
	}
    for _,e := range ret{
        l.list.Remove(e)
    }
}

//Get get all id
func (l *IDList)Get()(ret []string){
    for e := l.list.Front(); e != nil; e = e.Next() {
        ret = append(ret,e.Value.(string))
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
