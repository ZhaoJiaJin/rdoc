package httpapi

import (
	"fmt"
    log "github.com/sirupsen/logrus"
    "go.etcd.io/etcd/raft/raftpb"
	"net/http"
    "time"
	"rdoc/db"
)

var (
	HttpDB *db.DB // HTTP API endpoints operate on this database
    confChangeC chan<- raftpb.ConfChange
)

// Require Store form parameter value of specified key to *val and return true; if key does not exist, set HTTP status 400 and return false.
func Require(w http.ResponseWriter, r *http.Request, key string, val *string) bool {
	*val = r.FormValue(key)
	if *val == "" {
		http.Error(w, fmt.Sprintf("Please pass POST/PUT/GET parameter value of '%s'.", key), 400)
		return false
	}
	return true
}

// ServeHttpAPI HTTP server and block until the server shuts down. Panic on error.
func ServeHttpAPI(docdb *db.DB, port int, confChangeC chan<- raftpb.ConfChange, errorC <-chan error,authToken string) {
	var err error
	HttpDB = docdb

	var authWrap func(http.HandlerFunc) http.HandlerFunc
	if authToken != "" {
		log.Info("API endpoints now require the pre-shared token in Authorization header.")
		authWrap = func(originalHandler http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
                starttime := time.Now().UnixNano()
                log.Infof("receive new %v request %v",r.Method,r.URL.Path)
				if "token "+authToken != r.Header.Get("Authorization") {
					http.Error(w, "", http.StatusUnauthorized)
					return
				}
				originalHandler(w, r)
                endtime := time.Now().UnixNano()
                log.Infof("process request %v request %v cost:%v",r.Method,r.URL.Path,(endtime - starttime)/1e6)
			}
		}
	}else {
		log.Info("API endpoints do not require Authorization header.")
		authWrap = func(originalHandler http.HandlerFunc) http.HandlerFunc {
            return func(w http.ResponseWriter, r *http.Request) {
                starttime := time.Now().UnixNano()
                log.Infof("receive new %v request %v",r.Method,r.URL.Path)
				originalHandler(w, r)
                endtime := time.Now().UnixNano()
                log.Infof("process request %v request %v cost:%v",r.Method,r.URL.Path,(endtime - starttime)/1e6)
			}

			//return originalHandler
		}
	}
	// collection management (stop-the-world)
	http.HandleFunc("/create", authWrap(Create))
	http.HandleFunc("/rename", authWrap(Rename))
	http.HandleFunc("/drop", authWrap(Drop))
	http.HandleFunc("/all", authWrap(All))
	// query
	http.HandleFunc("/query", authWrap(Query))
	http.HandleFunc("/count", authWrap(Count))
	// document management
	http.HandleFunc("/insert", authWrap(Insert))
	http.HandleFunc("/get", authWrap(Get))
    //TODO
	//http.HandleFunc("/getpage", authWrap(GetPage))
	http.HandleFunc("/update", authWrap(Update))
	http.HandleFunc("/merge", authWrap(Merge))
	http.HandleFunc("/delete", authWrap(Delete))
	// index management (stop-the-world)
	http.HandleFunc("/index", authWrap(Index))
	http.HandleFunc("/indexes", authWrap(Indexes))
	http.HandleFunc("/unindex", authWrap(Unindex))
	//http.HandleFunc("/shutdown", authWrap(Shutdown))

    //TODO:raft api
    http.HandleFunc("/addnode", authWrap(AddNode))
    http.HandleFunc("/delnode", authWrap(DelNode))

	iface := "all interfaces"
    bind := ""

	log.Infof("Will listen on %s (HTTP), port %d.", iface, port)
    go func(){
        if err = http.ListenAndServe(fmt.Sprintf("%s:%d", bind, port), nil); err != nil{
            log.Fatal(err)
        }
    }()
	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}

