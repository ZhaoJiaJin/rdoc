package db

import(
    "testing"
    "rdoc/utils"
)


func TestIndexDoc(t *testing.T){
    path := "a,b,c"
    idx := NewIndex(path)
    id := RandID()

    data := `{"a": {"b":{"c":11}}, "a1": 2}`
	d, err := NewDoc([]byte(data))
    if err != nil{
        t.Fatal(err)
    }
    t.Log("add doc index id:",id)
    idx.IndexDoc(id,d)
    val := utils.StrHash("11")
    ids,err := idx.Query(val,10)
    if err != nil{
        t.Fatal(err)
    }
    if len(ids) != 1 || ids[0] != id{
        t.Fatal("index query failed")
    }

    idx.UnIndex(id)
    ids,err = idx.Query(val,10)
    if len(ids) != 0{
        t.Fatal("index query expect 0")
    }
}
