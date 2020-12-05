package avad

import (
	"ava/core"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/orcaman/concurrent-map"
	"github.com/phuslu/log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//ip--conn对应map
var allconn = make(map[string]*websocket.Conn)
var wsStatus = cmap.New()

func DialWs(addr string) {
	host:=strings.Split(addr,":")[0]


	wsStatus.Set(host, false)
	for {
		if status, ok := wsStatus.Get(host); ok {
			tmp := status.(bool)
			if !tmp {
				u := url.URL{Scheme: "ws", Host: addr, Path: "/echo"}
				c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
				if err != nil {
					log.Printf("尝试连接节点ws通道%s失败,10s后重试:\n", addr)
					wsStatus.Set(host, false)
					time.Sleep(10 * time.Second)
					continue
				}
				host:=strings.Split(u.Host,":")[0]
				wsStatus.Set(host, true)
				allconn[host] = c
				log.Debug().Msgf("已创建连接节点ws通道%s\n", addr)
			}
		}

	}
}

func DialS5(listenTarget string) {
	var RemoteConn net.Conn
	var err error
	for {
		for {
			RemoteConn, err = net.Dial("tcp", listenTarget)
			if err == nil {
				break
			}
		}
		go Handshake(RemoteConn)
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
			log.Debug().Msgf("节点消息投送失败,触发重新连接节点")
			msg = "投送失败,触发重新连接节点"
			wsStatus.Set(addr, false)
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

func staus(w http.ResponseWriter, r *http.Request)  {
	rv:=make(map[string]string)
	for k, v := range allconn {
		rv[k]=v.LocalAddr().String()
	}


	err := json.NewEncoder(w).Encode(rv)
	if err != nil {
		//... handle error
		panic(err)
	}

}



func staus2(w http.ResponseWriter, r *http.Request)  {
	err := json.NewEncoder(w).Encode(wsStatus)
	if err != nil {
		//... handle error
		panic(err)
	}

}


func DLocal(addrs []string) {
	//连接websocket

	for _, host := range addrs {
		addr := strings.Join([]string{host, ":", core.WsPort}, "")
		go DialWs(addr)
	}

	//连接内网穿透
	for _, host := range addrs {
		addr := strings.Join([]string{host, ":", core.TcpPort}, "")
		go DialS5(addr)
	}

	http.HandleFunc("/exectask", handel)
	http.HandleFunc("/status", staus)
	http.HandleFunc("/status2", staus2)
	addr := strings.Join([]string{"localhost", ":", core.Web}, "")
	http.ListenAndServe(addr, nil)

}
