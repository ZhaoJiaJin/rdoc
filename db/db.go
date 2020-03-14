package db

import (
	"sync"
    "go.etcd.io/etcd/etcdserver/api/snap"
    "encoding/json"
)

// DB database instance
type DB struct {
    proposeC chan<- string
    snapshotter *snap.Snapshotter
	//cols sync.Map
    cols map[string]*Col
    resCMap map[string]chan *OpeRet
    reslock sync.RWMutex
    sync.RWMutex
}

// NewDB create new database instance
func NewDB()*DB{
    return &DB{
        cols:make(map[string]*Col),
    }
}


func (db *DB)MarshalJSON() ([]byte, error) {
    db.RLock()
    defer db.RUnlock()
    return json.Marshal(db.cols)
}

func (db *DB)UnmarshalJSON(b []byte) error {
    db.Lock()
    defer db.Unlock()
    if err := json.Unmarshal(b,&(db.cols));err != nil{
        return err
    }
    return nil
}

func (db *DB) loadOrStore(key string, value *Col) (col *Col,loaded bool) {
    //db.Lock()
    //defer db.Unlock()
    col,loaded = db.cols[key]
    if loaded{
        return
    }
    db.cols[key] = value
    return value, false
}

func (db *DB) load(key string) (col *Col,ok bool) {
    //db.RLock()
    //defer db.RUnlock()
    col,ok = db.cols[key]
    return
}

func (db *DB) store(key string, value *Col) {
    //db.Lock()
    //defer db.Unlock()
	db.cols[key] = value
}

func (db *DB) del(key string) {
    //db.Lock()
    //defer db.Unlock()
    delete(db.cols,key)
}

//CreateCol create collection
func (db *DB) CreateCol(colname string) error {
	_, loaded := db.loadOrStore(colname, NewCol())
	if loaded {
		return ErrColExist
	}
    db.store(colname,NewCol())
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
	/*_, ok = db.load(newname)
	if ok {
		return ErrColExist
	}*/
    _,loaded := db.loadOrStore(newname, col)
    if loaded{
		return ErrColExist
    }
	db.del(oldname)
	return nil
}

//GetAllCol get names of all collections
func (db *DB) GetAllCol() []string {
	res := make([]string, 0)
    db.RLock()
    defer db.RUnlock()
    for key := range db.cols{
		res = append(res, key)
    }
	return res
}

//InsertDoc insert doc into a collection
//TODO: deal with id in distributed system
func (db *DB) InsertDoc(colname string, data []byte,id string) (string, error) {
	col, ok := db.load(colname)
	if !ok {
		return "", ErrColNotExist
	}
	return col.AddDoc(id,data)
}

//UpdateDoc update doc
func (db *DB) UpdateDoc(colname string, data []byte, ids string) error {
	col, ok := db.load(colname)
	if !ok {
		return ErrColNotExist
	}
	return col.UpdateDoc(ids,data)
}

//MergeDoc merge document with given document
func (db *DB) MergeDoc(colname string, data []byte, ids string) error {
	col, ok := db.load(colname)
	if !ok {
		return ErrColNotExist
	}
	return col.MergeDoc(ids,data)
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
    db.RLock()
    defer db.RUnlock()
	col, ok := db.load(colname)
	if !ok {
		return 0, ErrColNotExist
	}
	ids, err := col.QueryDocID(data)
	return len(ids), err
}

//QueryDoc query doc
func (db *DB) QueryDocByID(colname string, id string) (map[string]interface{}, error) {
    db.RLock()
    defer db.RUnlock()
	var res map[string]interface{}
	col, ok := db.load(colname)
	if !ok {
		return res, ErrColNotExist
	}
	res = make(map[string]interface{})
    if cl :=  col.ReadDoc(id); cl != nil{
	    res[id] = cl.Data
    }
	return res, nil
}



//QueryDoc query doc
func (db *DB) QueryDoc(colname string, data []byte) (map[string]interface{}, error) {
    db.RLock()
    defer db.RUnlock()
	var res map[string]interface{}
	col, ok := db.load(colname)
	if !ok {
		return res, ErrColNotExist
	}
	ids, err := col.QueryDocID(data)
	if err != nil {
		return res, err
	}
	res = make(map[string]interface{})
	for id := range ids {
        if tmpcl := col.ReadDoc(id);tmpcl != nil{
		    res[id] = tmpcl.Data
        }
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
    db.RLock()
    defer db.RUnlock()
	res := make([]string, 0)
	col, ok := db.load(colname)
	if !ok {
		return res, ErrColNotExist
	}
	return col.GetAllIndex()
}
