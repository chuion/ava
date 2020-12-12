package avad

import (
	"encoding/json"
	cmap "github.com/orcaman/concurrent-map"
	"net/http"
)

func webWsStatus(w http.ResponseWriter, r *http.Request) {
	rv := map[string]cmap.ConcurrentMap{}
	rv["ws"] = wsStatus
	rv["tcp"] = tcpStatus
	err := json.NewEncoder(w).Encode(rv)
	if err != nil {
		//... handle error
		panic(err)
	}

}

func webWorkerMapR(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(workerMapR)
	if err != nil {
		//... handle error
		panic(err)
	}

}

func info(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(Ver)
	if err != nil {
		//... handle error
		panic(err)
	}

}
