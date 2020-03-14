// Document management handlers.

package httpapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
    "rdoc/db"
)

// Insert a document into collection.
func Insert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
	var col, doc string
	if !Require(w, r, "col", &col) {
		return
	}
	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	doc = string(bodyBytes)
	if doc == "" && !Require(w, r, "doc", &doc) {
		return
	}
	/*id, err := HttpDB.InsertDoc(col,[]byte(doc))
	if err != nil {
		http.Error(w, fmt.Sprint(err), 400)
		return
	}*/

    ope := db.Operate{OpeType:db.INSERTDOC, ColName:col,Data:[]byte(doc)}
    reschan, _ := HttpDB.Propose(ope)
    var res *db.OpeRet
    for res = range reschan{
    }
    if res.Err != nil{
        http.Error(w, res.Err.Error(), http.StatusInternalServerError)
    }else{
	    w.WriteHeader(201)
        w.Write([]byte(res.Msg))
    }
}

// Get find and retrieve a document by ID.
func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
	var col, id string
	if !Require(w, r, "col", &col) {
		return
	}
	if !Require(w, r, "id", &id) {
		return
	}
	doc, err := HttpDB.QueryDocByID(col,id)
	if doc == nil {
		http.Error(w, fmt.Sprintf("No such document ID %d.", id), 404)
		return
	}
	resp, err := json.Marshal(doc)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
		return
	}
	w.Write(resp)
}

// Update a documents.
func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
	var col, id, doc string
	if !Require(w, r, "col", &col) {
		return
	}
	if !Require(w, r, "id", &id) {
		return
	}
	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	doc = string(bodyBytes)
	if doc == "" && !Require(w, r, "doc", &doc) {
		return
	}
    /*err := HttpDB.UpdateDoc(col, []byte(doc), id)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
		return
	}*/

    ope := db.Operate{OpeType:db.UPDATEDOC, ColName:col,Data:[]byte(doc),IDs:id}
    reschan, _ := HttpDB.Propose(ope)
    var res *db.OpeRet
    for res = range reschan{
    }
	if res.Err != nil {
		http.Error(w, fmt.Sprint(res.Err), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}


// Merge merge update documents.
func Merge (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
	var col, id, doc string
	if !Require(w, r, "col", &col) {
		return
	}
	if !Require(w, r, "id", &id) {
		return
	}
	defer r.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	doc = string(bodyBytes)
	if doc == "" && !Require(w, r, "doc", &doc) {
		return
	}
    /*err := HttpDB.MergeDoc(col, []byte(doc), id)
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
		return
	}*/
    ope := db.Operate{OpeType:db.MERGEDOC, ColName:col,Data:[]byte(doc),IDs:id}
    reschan, _ := HttpDB.Propose(ope)
    var res *db.OpeRet
    for res = range reschan{
    }
	if res.Err != nil {
		http.Error(w, fmt.Sprint(res.Err), http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// Delete a document.
func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
	var col, id string
	if !Require(w, r, "col", &col) {
		return
	}
	if !Require(w, r, "id", &id) {
		return
	}
	//HttpDB.DeleteDoc(col,id)
    ope := db.Operate{OpeType:db.DELETEDOC, ColName:col,IDs:id}
    reschan, _ := HttpDB.Propose(ope)

    var res *db.OpeRet
    for res = range reschan{
    }
	if res.Err != nil {
		http.Error(w, fmt.Sprint(res.Err), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
    }
}

