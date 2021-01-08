package avad

import (
	"ava/core"
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
	err := json.NewEncoder(w).Encode(AllInfo)
	if err != nil {
		//... handle error
		panic(err)
	}

}

type webInfo struct {
	Host string `json:"Host"`
	Business []string `json:"Business"`
	Status   bool     `json:"Status"`
	core.PcInfo
}

func getAllInfo(w http.ResponseWriter, r *http.Request) {
	rv := []webInfo{}

	ch := wsStatus.IterBuffered()
	for item := range ch {
		host := item.Key
		tmp, _ := wsStatus.Get(host)
		sta := tmp.(bool)
		infoOne := webInfo{
			Host: host,
			Business: workerMapR[host],
			Status:   sta,
			PcInfo:   AllInfo[host],
		}
		rv=append(rv,infoOne)
	}
	err := json.NewEncoder(w).Encode(rv)
	if err != nil {
		//... handle error
		panic(err)
	}
}
