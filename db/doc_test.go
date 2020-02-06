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
	res := GetIn(d.data, []string{"a", "hah"})
	t.Log(res)
}
