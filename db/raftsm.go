package db

import(
    "go.etcd.io/etcd/etcdserver/api/snap"
    "encoding/gob"
	"bytes"
    "encoding/json"
    log "github.com/sirupsen/logrus"
)

const(
    CREATECOL = iota
    REMOVECOL
    RENAMECOL
    INSERTDOC
    UPDATEDOC
    MERGEDOC
    DELETEDOC
    CREATEIDX
    RMIDX
)

//Operate database operate request
type Operate struct{
    OID string
    OpeType int
    ColName string
    Data []byte
    IDs string
    Path string
}

//OpeRet operate result
type OpeRet struct{
    Msg string
    Err error
}

//NewDBWithRaft create db instance with raft support
func NewDBWithRaft(snapshotter *snap.Snapshotter, proposeC chan<- string, commitC <-chan *string, errorC <-chan error)*DB{
    db := NewDB()
    db.resCMap = make(map[string]chan *OpeRet)
    db.proposeC = proposeC
    db.snapshotter = snapshotter

    db.readCommits(commitC,errorC)
    go db.readCommits(commitC,errorC)
    return db
}

func (db *DB) readCommits(commitC <-chan *string, errorC <-chan error){
    for data := range commitC{
        if data == nil{
            snapshot,err := db.snapshotter.Load()
            if err == snap.ErrNoSnapshot{
                return
            }
            if err != nil{
                log.Panic(err)
            }
            log.Infof("loading snapshot at term %d and index %d",snapshot.Metadata.Term, snapshot.Metadata.Index)
            if err := db.recoverFromSnapshot(snapshot.Data); err != nil{
                log.Panic(err)
            }
            return
        }

        var dataOpe Operate
		dec := gob.NewDecoder(bytes.NewBufferString(*data))
        if err := dec.Decode(&dataOpe); err != nil{
			log.Fatalf("readCommits: could not decode message (%v)", err)
        }
        res := db.applyOpe(&dataOpe)
        db.reslock.RLock()
        reschan,ok := db.resCMap[dataOpe.OID]
        db.reslock.RUnlock()
        if ok{
            reschan <- res
            close(reschan)
            db.reslock.Lock()
            delete(db.resCMap,dataOpe.OID)
            db.reslock.Unlock()
        }
    }
    if err, ok := <-errorC; ok{
        log.Fatal("errorC:",err)
    }
}

func (db *DB)applyOpe(ope *Operate)(res *OpeRet){
    db.Lock()
    defer db.Unlock()
    res = &OpeRet{}
    switch ope.OpeType{
    case CREATECOL:
        res.Err = db.CreateCol(ope.ColName)
    case REMOVECOL:
        db.RemoveCol(ope.ColName)
    case RENAMECOL:
        res.Err = db.RenameCol(ope.ColName,string(ope.Data))
    case INSERTDOC:
        res.Msg, res.Err = db.InsertDoc(ope.ColName, ope.Data,ope.OID)
    case UPDATEDOC:
        res.Err = db.UpdateDoc(ope.ColName, ope.Data, ope.IDs)
    case MERGEDOC:
        res.Err = db.MergeDoc(ope.ColName, ope.Data, ope.IDs)
    case DELETEDOC:
        res.Err = db.DeleteDoc(ope.ColName,ope.IDs)
    case CREATEIDX:
        res.Err = db.CreateIndex(ope.ColName,ope.Path)
    case RMIDX:
        res.Err = db.RemoveIndex(ope.ColName,ope.Path)
    default:
        res.Err = ErrUnSuppOpe
    }
    return
}

// Propose propose a new request
func (db *DB) Propose(ope Operate)(chan *OpeRet, string){
    //if ope.OID == ""{
    ope.OID = RandID()
    //}
    rchan := make(chan *OpeRet,1)
    db.reslock.Lock()
    db.resCMap[ope.OID] = rchan
    db.reslock.Unlock()

    var buf bytes.Buffer
    if err := gob.NewEncoder(&buf).Encode(ope); err != nil{
        log.Fatal("Propose:",err)
    }
    db.proposeC <- buf.String()
    return rchan, ope.OID
}

func (db *DB) GetSnapshot()([]byte,error){
    return json.Marshal(db)
}


func (db *DB) recoverFromSnapshot(snapshot []byte)(error){
    return json.Unmarshal(snapshot,&db)
}
