package avad

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/orcaman/concurrent-map"
	"log"
	"net/http"
	"net/url"
	"time"
)

var addrs = []string{"localhost:8080", "localhost:8081"}
var allconns = make(map[string]*websocket.Conn)
var connsStatus = cmap.New()


func DialOne(addr string) {
	connsStatus.Set(addr, false)
	for {
		if status, ok := connsStatus.Get(addr); ok {
			tmp := status.(bool)
			if !tmp {
				u := url.URL{Scheme: "ws", Host: addr, Path: "/echo"}
				c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
				if err != nil {
					log.Printf("尝试连接节点%s失败,10s后重试:\n", addr)
					connsStatus.Set(addr, false)
					time.Sleep(10 * time.Second)
					continue
				}
				connsStatus.Set(addr, true)
				allconns[addr] = c
				fmt.Printf("已创建连接节点%s\n", addr)
			}
		}

	}
}

func DLocal() {
	for _, addr := range addrs {
		go DialOne(addr)
	}

	go Socks5d()


	http.HandleFunc("/start", handel)
	addr := "localhost:4000"
	log.Fatal(http.ListenAndServe(addr, nil))

}

type Task struct {
	Cmd  string
	Args string
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
	addr := p.Cmd
	c := allconns[addr]

	err = c.WriteMessage(websocket.TextMessage, []byte("这是一条测试"))
	msg := "投送成功"
	if err != nil {
		fmt.Println("节点消息投送失败,触发重新连接节点")
		//panic(err)
		msg = "投送失败,触发重新连接节点"
		connsStatus.Set(addr, false)

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
