package db

import(
    "testing"
    //"encoding/json"
)

/*func createDoc(t *testing.T)[]byte{
    data := make(map[string]interface{})
}*/

func TestDB(t *testing.T){
    db := NewDB()
    colname := "test-col"
    //test create
    err := db.CreateCol(colname)
    if err != nil{
        t.Fatal(err)
    }

    //test create index
    err = db.CreateIndex(colname,"")
    if err != nil{
        t.Fatal(err)
    }


    data := `{"a": {"b":1}, "a1": 2}`
    id,err := db.InsertDoc(colname,[]byte(data))
    if err != nil{
        t.Fatal(err)
    }
    t.Log("add doc:",id)
    data = `{"a": {"b":1}, "a1": 2}`
    id,err = db.InsertDoc(colname,[]byte(data))
    if err != nil{
        t.Fatal(err)
    }
    t.Log("add doc:",id)
}
