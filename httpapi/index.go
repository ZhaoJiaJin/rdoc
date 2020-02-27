// Index management handlers.

package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
    "rdoc/db"
)

// Index Put an index on a document path.
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
	var col, path string
	if !Require(w, r, "col", &col) {
		return
	}
	if !Require(w, r, "path", &path) {
		return
	}
	/*if err := HttpDB.CreateIndex(col, path); err != nil {
		http.Error(w, fmt.Sprint(err), 400)
		return
	}*/

    ope := db.Operate{OpeType:db.CREATEIDX, ColName:col,Path:path}
    HttpDB.Propose(ope)
    //TODO: return result
	w.WriteHeader(201)
}

// Return all indexed paths.
func Indexes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
	var col string
	if !Require(w, r, "col", &col) {
		return
	}
    indexes,err := HttpDB.GetAllIndex(col)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	resp, err := json.Marshal(indexes)
	if err != nil {
		http.Error(w, fmt.Sprint("Server error."), 500)
		return
	}
	w.Write(resp)
}

// Unindex Remove an indexed path.
func Unindex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
	var col, path string
	if !Require(w, r, "col", &col) {
		return
	}
	if !Require(w, r, "path", &path) {
		return
	}
	/*if err := HttpDB.RemoveIndex(col,path); err != nil {
		http.Error(w, fmt.Sprint(err), 400)
		return
	}*/

    ope := db.Operate{OpeType:db.RMIDX, ColName:col,Path:path}
    HttpDB.Propose(ope)
    //TODO: return result
}
