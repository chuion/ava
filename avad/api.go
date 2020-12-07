package avad

import (
	"encoding/json"
	"net/http"
)

func webWsConns(w http.ResponseWriter, r *http.Request) {
	rv := make(map[string]string)
	for k, v := range wsConns {
		rv[k] = v.LocalAddr().String()
	}

	err := json.NewEncoder(w).Encode(rv)
	if err != nil {
		//... handle error
		panic(err)
	}

}

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
