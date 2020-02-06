package db

import(
	"encoding/json"
)


type Doc struct{
    data map[string]interface{}
}



func NewDoc(data []byte)(*Doc,error){
	d := Doc{
        data:make(map[string]interface{}),
    }
	err := json.Unmarshal(data,&d.data)
	return &d, err
}

func GetIn(doc interface{},path []string)(ret []interface{}){
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
