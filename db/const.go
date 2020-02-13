package db

import (
	"errors"
)

var (
	//ErrColExist Collection already exist
	ErrColExist = errors.New("Collection exist")
	//ErrColNotExist Collection doesn't exist
	ErrColNotExist = errors.New("Collection doesn't exist")


	//ErrIDXExist index already exist
	ErrIDXExist = errors.New("index exist")
    //ErrNotIDX query path is not indexed
    ErrNotIDX = errors.New("path is not indexed")
    //ErrLimit wrong limit number
    ErrLimit = errors.New("wrong limit number")
    ErrLimitType = errors.New("wrong limit type")
    ErrNoSubQuery = errors.New("expecting subquery")
    ErrRangeType = errors.New("range type wrong")
    ErrRangeMiss = errors.New("miss range values from-to or int-to")
)

const(
    INDEX_PATH_SEP = ","
)
