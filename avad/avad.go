package avad

import (
	"ava/core"
	"github.com/orcaman/concurrent-map"
	"net/http"
	"strings"
)

//ip--conn对应map
var wsConns = cmap.New()
var wsStatus = cmap.New()
var tcpStatus = cmap.New()

func Manger(addrs []string) {

	for _, host := range addrs {
		tcpStatus.Set(host, false)
		wsStatus.Set(host, false)
		wsConns.Set(host, nil)
	}

	go ping()

	http.HandleFunc("/exectask", taskRouter)
	http.HandleFunc("/webWsStatus", webWsStatus)
	http.HandleFunc("/webWorkerMapR", webWorkerMapR)
	http.HandleFunc("/info", info)
	http.HandleFunc("/v1/allInfo", getAllInfo)
	http.HandleFunc("/v1/proxy", getProxyInfo)
	http.Handle("/", http.FileServer(http.Dir("dist")))

	addr := strings.Join([]string{"0.0.0.0", ":", core.Web}, "")
	http.ListenAndServe(addr, nil)

}
