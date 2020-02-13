package db

import(
    "testing"
    "rdoc/utils"
)

func TestDoc(t *testing.T){
    col := NewCol()

    data := `{"a": {"b":1}, "a1": 2}`
    id1, err := col.AddDoc([]byte(data))
    if err != nil{
        t.Fatal(err)
    }
    data1 := `{"a": {"b":1}, "a1": 1000}`
    id2, err := col.AddDoc([]byte(data1))
    if err != nil{
        t.Fatal(err)
    }
 
    lookupval := utils.StrHash("1")
    ids,err := col.Query("a,b",lookupval,10)
    if err != ErrNotIDX{
        t.Fatal("expect failed operation")
    }

    err = col.CreateIndex("a,b")
    if err != nil{
        t.Fatal(err)
    }
    ids,err = col.Query("a,b",lookupval,10)
    if err != nil{
        t.Fatal(err)
    }

    t.Log(ids,id1,id2)
    if len(ids) != 2 {
        t.Fatal("wrong query result")
    }

    d1 := col.ReadDoc(id2)
    t.Log(d1)

}
