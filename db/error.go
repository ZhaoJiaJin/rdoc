package db

import (
	"errors"
)

var (
	//ErrColExist Collection already exist
	ErrColExist = errors.New("Collection exist")
	//ErrColNotExist Collection doesn't exist
	ErrColNotExist = errors.New("Collection doesn't exist")
)
