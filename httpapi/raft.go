package httpapi


import(
    "net/http"
    "strconv"
    "go.etcd.io/etcd/raft/raftpb"
    log "github.com/sirupsen/logrus"
)

//AddNode add node to raft cluster
func AddNode(w http.ResponseWriter, r *http.Request) {
	var nodeId, url string
    if ok := Require(w,r,"nodeID",&nodeId); !ok{
        log.Warnf("AddNode:fail to read nodeID")
        http.Error(w, "no nodeID", http.StatusBadRequest)
        return
    }
    if ok := Require(w,r,"url",&url); !ok{
        log.Warnf("AddNode:fail to read url")
        http.Error(w, "no url", http.StatusBadRequest)
        return
    }
    nodeIdint,err := strconv.ParseUint(nodeId,10,64)
    if err != nil{
        log.Warnf("AddNode:wrong nodeID type")
        http.Error(w, "wrong nodeID type", http.StatusBadRequest)
        return
    }

    log.Infof("httpapi AddNode:%v %v",nodeIdint, url)
	cc := raftpb.ConfChange{
		Type:    raftpb.ConfChangeAddNode,
		NodeID:  nodeIdint,
		Context: []byte(url),
	}
	confChangeC <- cc

	// As above, optimistic that raft will apply the conf change
	w.WriteHeader(http.StatusNoContent)
}


func DelNode(w http.ResponseWriter, r *http.Request) {
	var nodeId string
    if ok := Require(w,r,"nodeID",&nodeId); !ok{
        log.Warnf("DelNode:fail to read nodeID")
        http.Error(w, "no nodeID", http.StatusBadRequest)
        return
    }
    nodeIdint,err := strconv.ParseUint(nodeId,10,64)
    if err != nil{
        log.Warnf("DelNode:wrong nodeID type")
        http.Error(w, "wrong nodeID type", http.StatusBadRequest)
        return
    }


	cc := raftpb.ConfChange{
		Type:   raftpb.ConfChangeRemoveNode,
		NodeID: nodeIdint,
	}
	confChangeC <- cc

	// As above, optimistic that raft will apply the conf change
	w.WriteHeader(http.StatusNoContent)
}
