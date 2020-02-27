package db

import (
	"encoding/json"
)

type Doc struct {
	Data map[string]interface{}
}

func NewDoc(data []byte) (*Doc, error) {
	d := Doc{
		Data: make(map[string]interface{}),
	}
	err := json.Unmarshal(data, &d.Data)
	return &d, err
}

func GetIn(doc interface{}, path []string) (ret []interface{}) {
	docMap, ok := doc.(map[string]interface{})
	if !ok {
		return
	}
	var thing interface{} = docMap
	// Get into each path segment
	for i, seg := range path {
		if aMap, ok := thing.(map[string]interface{}); ok {
			thing = aMap[seg]
		} else if anArray, ok := thing.([]interface{}); ok {
			for _, element := range anArray {
				ret = append(ret, GetIn(element, path[i:])...)
			}
			return ret
		} else {
			return nil
		}
	}
	switch thing := thing.(type) {
	case []interface{}:
		return append(ret, thing...)
	default:
		return append(ret, thing)
	}
}



//Merge merge two documents
func (d *Doc)Merge(ndoc *Doc){
    merge(d.Data, ndoc.Data)
}

func merge(d1 map[string]interface{},d2 map[string]interface{}){
    /*d1map, ok := d1.(map[string]interface{})
    if !ok{
        return
    }
    d2map, ok := d2.(map[string]interface{})
    if !ok{
        return
    }*/
    for k,v := range d2{
        d1v,ok := d1[k]
        if !ok{
            d1[k] = v
        }else{
            d1vmap,d1vok := d1v.(map[string]interface{})
            d2vmap,d2vok := d2[k].(map[string]interface{})
            if d1vok && d2vok{
                merge(d1vmap,d2vmap)
            }else{
                d1[k] = v
            }
        }
    }
}
