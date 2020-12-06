package avad

import (
	"ava/core"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/orcaman/concurrent-map"
	"github.com/phuslu/log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//ip--conn对应map
var allconn = make(map[string]*websocket.Conn)
var nodeTask = make(map[string]core.LauncherConf)
var wsStatus = cmap.New()
var tcpStatus = cmap.New()

func DialWs(addr string) {
	host := strings.Split(addr, ":")[0]

	wsStatus.Set(host, false)
	for {
		status, _ := wsStatus.Get(host)
		if !status.(bool) {
			u := url.URL{Scheme: "ws", Host: addr, Path: "/echo"}
			c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Debug().Msgf("尝试连接节点ws通道%s失败,%s后重试:\n", addr, pongWait)
				time.Sleep(pongWait)
				continue
			}
			host := strings.Split(u.Host, ":")[0]
			wsStatus.Set(host, true)
			allconn[host] = c
			log.Debug().Msgf("已创建连接节点ws通道%s\n", addr)

			p := core.LauncherConf{}
			err = c.ReadJSON(&p)
			if err != nil {
				log.Debug().Msgf("接收节点: %s注册信息失败", host)
			}
			log.Debug().Msgf("接收节点: %s注册信息成功,可运行%s", host, p.Worker)
			nodeTask[host] = p
		}
	}

}

type Task struct {
	Route string
	Cmd   string
	Args  string
}

type rusult struct {
	Code int
	Msg  string
}

func handel(w http.ResponseWriter, r *http.Request) {
	var p Task
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//todo 往哪个连接发应解析业务逻辑
	addr := p.Route
	msg := "投送成功"
	if c, ok := allconn[addr]; ok {
		err = c.WriteJSON(p)
		if err != nil {
			log.Debug().Msgf("投送失败,节点可能已不在线")
			msg = "投送失败,节点可能已不在线"

		}
	} else {
		msg = fmt.Sprintf("未找到%s对应的socket连接", p.Route)
	}

	rv := rusult{
		Code: 200,
		Msg:  msg,
	}
	err = json.NewEncoder(w).Encode(rv)
	if err != nil {
		//... handle error
		panic(err)
	}

}

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

func DLocal(addrs []string) {
	//连接websocket

	for _, host := range addrs {
		tcpStatus.Set(host, false)
		wsStatus.Set(host, false)
	}

	go ping()

	for _, host := range addrs {
		addr := strings.Join([]string{host, ":", core.WsPort}, "")
		go DialWs(addr)
	}

	//连接内网穿透
	for _, host := range addrs {
		addr := strings.Join([]string{host, ":", core.TcpPort}, "")
		go connectForSocks(addr)
	}

	http.HandleFunc("/exectask", handel)
	http.HandleFunc("/webWsConns", webWsConns)
	http.HandleFunc("/webWsStatus", webWsStatus)
	http.HandleFunc("/webNodeTask", webNodeTask)
	http.HandleFunc("/webTcpStatus", webTcpStatus)
	addr := strings.Join([]string{"localhost", ":", core.Web}, "")
	http.ListenAndServe(addr, nil)

}
