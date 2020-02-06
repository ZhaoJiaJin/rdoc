package db

import (
	"sync"
)

// DB database instance
type DB struct {
	cols sync.Map
}

func (db *DB) loadOrStore(key string, value *Col) (*Col, bool) {
	res, loaded := db.cols.LoadOrStore(key, value)
	return res.(*Col), loaded
}

func (db *DB) load(key string) (*Col, bool) {
	res, ok := db.cols.Load(key)
	if !ok {
		return nil, ok
	}
	return res.(*Col), ok
}

func (db *DB) store(key string, value *Col) {
	db.cols.Store(key, value)
}

func (db *DB) del(key string) {
	db.cols.Delete(key)
}

//CreateCol create collection
func (db *DB) CreateCol(colname string) error {
	_, loaded := db.loadOrStore(colname, NewCol())
	if loaded {
		return ErrColExist
	}
	return nil
}

//RemoveCol remove collection
func (db *DB) RemoveCol(colname string) {
	db.del(colname)
}

//RenameCol rename a collection
func (db *DB) RenameCol(oldname, newname string) error {
	col, ok := db.load(oldname)
	if !ok {
		return ErrColNotExist
	}
	_, ok = db.load(newname)
	if ok {
		return ErrColExist
	}
	db.loadOrStore(newname, col)
	db.del(oldname)
	return nil
}

//GetAllCol get names of all collections
func (db *DB) GetAllCol() []string {
	res := make([]string, 0)
	db.cols.Range(func(key, value interface{}) bool {
		res = append(res, key.(string))
		return true
	})
	return res
}

//InsertDoc insert doc into a collection
func (db *DB) InsertDoc(colname string, data []byte) (string, error) {
	col, ok := db.load(colname)
	if !ok {
		return "", ErrColNotExist
	}
	return col.AddDoc(data)
}

//UpdateDoc update doc
func (db *DB) UpdateDoc(colname string, data []byte, ids string) error {
	col, ok := db.load(colname)
	if !ok {
		return ErrColNotExist
	}
	return col.UpdateDoc(data)
}

//MergeDoc merge document with given document
func (db *DB) MergeDoc(colname string, data []byte, ids string) error {
	col, ok := db.load(colname)
	if !ok {
		return ErrColNotExist
	}
	return col.MergeDoc(data)
}

//DeleteDoc insert doc into a collection
func (db *DB) DeleteDoc(colname string, ids string) error {
	col, ok := db.load(colname)
	if !ok {
		return ErrColNotExist
	}
	return col.DeleteDoc(ids)
}

//CountDoc query doc
func (db *DB) CountDoc(colname string, data []byte) (int, error) {
	col, ok := db.load(colname)
	if !ok {
		return 0, ErrColNotExist
	}
	ids, err := col.QueryDocID(data)
	return len(ids), err
}

//QueryDoc query doc
func (db *DB) QueryDoc(colname string, data []byte) (map[string]*Doc, error) {
	var res map[string]*Doc
	col, ok := db.load(colname)
	if !ok {
		return res, ErrColNotExist
	}
	ids, err := col.QueryDocID(data)
	if err != nil {
		return res, err
	}
	res = make(map[string]*Doc)
	for _, id := range ids {
		res[id] = col.ReadDoc(id)
	}
	return res, nil
}

//CreateIndex create index
func (db *DB) CreateIndex(colname string, path string) error {
	col, ok := db.load(colname)
	if !ok {
		return ErrColNotExist
	}
	return col.CreateIndex(path)
}

//RemoveIndex remove index
func (db *DB) RemoveIndex(colname, path string) error {
	col, ok := db.load(colname)
	if !ok {
		return ErrColNotExist
	}
	return col.RmIndex(path)
}

//GetAllIndex return all index names
func (db *DB) GetAllIndex(colname string) ([]string, error) {
	res := make([]string, 0)
	col, ok := db.load(colname)
	if !ok {
		return res, ErrColNotExist
	}
	return col.GetAllIndex()
}
