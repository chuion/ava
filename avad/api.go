package avad

import (
	"encoding/json"
	"net/http"
)

func webWsStatus(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(wsStatus)
	if err != nil {
		//... handle error
		panic(err)
	}

}

func webWorkerMap(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(workerMap)
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

func webTcpStatus(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(tcpStatus)
	if err != nil {
		//... handle error
		panic(err)
	}

}
