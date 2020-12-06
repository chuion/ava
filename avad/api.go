package avad

import (
	"encoding/json"
	"net/http"
)

func webWsConns(w http.ResponseWriter, r *http.Request) {
	rv := make(map[string]string)
	for k, v := range allconn {
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

func webNodeTask(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(nodeTask)
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
