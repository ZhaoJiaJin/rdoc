package db

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
    "container/list"
)

var (
	node *snowflake.Node
)

type IDList struct{
    list *list.List
}

func NewIDList()*IDList{
    return &IDList{
        list:list.New(),
    }
}

func (l *IDList)ADD(id string)

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

//Randstring generate randomid
func RandID() string {
	return fmt.Sprintf("%X", node.Generate().Int64())
}
