package db

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node
)

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
