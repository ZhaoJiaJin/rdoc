package db

import(
    "sync"
)

// DB database instance
type DB struct{
    cols map[string]*Col
    sync.RWMutex
}

func (db *DB)CreateCol()(error){
    return nil
}

func (db *DB)RemoveCol()(error){
    return nil
}

func (db *DB)GetCol()(*Col){
    return nil
}


