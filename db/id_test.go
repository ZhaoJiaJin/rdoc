package db

import(
    "testing"
)

func TestIDList(t *testing.T){
    l := NewIDList()
    l.Add("1")
    l.Add("2")
    t.Log(l.Get())
    l.Remove("1")
    res := l.Get()
    if len(res) != 1 || res[0] != "2"{
        t.Fatal("id list failed")
    }
}
