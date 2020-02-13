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
    //err := db.CreateIndex(colname,"")



    // test get all

    //test rename

    //test remove
}
