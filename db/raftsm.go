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
    //OID string
    OpeType int
    ColName string
    Data []byte
    IDs string
    Path string
}

//NewDBWithRaft create db instance with raft support
func NewDBWithRaft(snapshotter *snap.Snapshotter, proposeC chan<- string, commitC <-chan *string, errorC <-chan error)*DB{
    db := NewDB()
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
        msg,err := db.applyOpe(&dataOpe)
        if err != nil{
            log.Errorf("applyOpe:%v",err)
        }else{
            log.Infof("applyOpe:%v",msg)
        }
        //TODO:pass msg & err to http api
    }
    if err, ok := <-errorC; ok{
        log.Fatal("errorC:",err)
    }
}

func (db *DB)applyOpe(ope *Operate)(interface{},error){
    ok := ""
    //fail := "fail"
    switch ope.OpeType{
    case CREATECOL:
        err := db.CreateCol(ope.ColName)
        return ok,err
    case REMOVECOL:
        db.RemoveCol(ope.ColName)
        return "",nil
    case RENAMECOL:
        err := db.RenameCol(ope.ColName,string(ope.Data))
        return "",err
    case INSERTDOC:
        return db.InsertDoc(ope.ColName, ope.Data)
    case UPDATEDOC:
        err := db.UpdateDoc(ope.ColName, ope.Data, ope.IDs)
        return "",err
    case MERGEDOC:
        err := db.MergeDoc(ope.ColName, ope.Data, ope.IDs)
        return "",err
    case DELETEDOC:
        err := db.DeleteDoc(ope.ColName,ope.IDs)
        return "",err
    case CREATEIDX:
        err := db.CreateIndex(ope.ColName,ope.Path)
        return "",err
    case RMIDX:
        err := db.RemoveIndex(ope.ColName,ope.Path)
        return "",err
    default:
        return "",ErrUnSuppOpe
    }
}

func (db *DB) Propose(ope Operate){
    /*if ope.OID == ""{
        ope.OID = RandID()
    }*/
    var buf bytes.Buffer
    if err := gob.NewEncoder(&buf).Encode(ope); err != nil{
        log.Fatal("Propose:",err)
    }
    db.proposeC <- buf.String()
}

func (db *DB) GetSnapshot()([]byte,error){
    return json.Marshal(db)
}


func (db *DB) recoverFromSnapshot(snapshot []byte)(error){
    return json.Unmarshal(snapshot,&db)
}
