package db

import (
	"testing"
)

func TestNewDoc(t *testing.T) {
	data := `{"a": {"hah":1}, "b": 2}`
	d, err := NewDoc([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(d)
}

func TestDocGetIn(t *testing.T) {
	data := `{"a": {"hah":1}, "b": 2}`
	d, err := NewDoc([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	res := GetIn(d.Data, []string{"a", "hah"})
	t.Log(res)
}

func TestMergeDoc(t *testing.T){
	data := `{"a": {"hah":1}, "b": 2}`
	d, err := NewDoc([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
    t.Log(d)
	data1 := `{"a": {"hah":"9"}}`
	d1, err := NewDoc([]byte(data1))
	if err != nil {
		t.Fatal(err)
	}
    t.Log(d1)

    d.Merge(d1)
    t.Log(d)
    res := GetIn(d.Data,[]string{"a", "hah"})
    if len(res) != 1 || res[0].(string) != "9"{
        t.Log("res",res)
        t.Fatal("merge doc failed")
    }
    t.Log(d)
    t.Log(d1)
}
