// Query handlers.

package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

)

// Query Execute a query and return documents from the result.
func Query(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
	var col, q string
	if !Require(w, r, "col", &col) {
		return
	}
	if !Require(w, r, "q", &q) {
		return
	}
    res, err := HttpDB.QueryDoc(col,[]byte(q))
    if err != nil{
		http.Error(w, fmt.Sprint(err), 400)
		return
    }
	resp, err := json.Marshal(res)
	if err != nil {
		http.Error(w, fmt.Sprintf("Server error: query returned invalid structure"), 500)
		return
	}
	w.Write([]byte(string(resp)))
}

// Count Execute a query and return number of documents from the result.
func Count(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
	var col, q string
	if !Require(w, r, "col", &col) {
		return
	}
	if !Require(w, r, "q", &q) {
		return
	}
	reslen, err := HttpDB.CountDoc(col, []byte(q))
    if err != nil{
		http.Error(w, fmt.Sprint(err), 400)
		return
    }
	w.Write([]byte(strconv.Itoa(reslen)))
}
