package db

import(
    "encoding/json"
)


type Doc map[string]interface{}


func NewDoc(data []byte)(*Doc,error){
    d := new(Doc)
    err := json.Unmarshal(data,d)
    return d, err
}
