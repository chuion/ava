package avad

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"time"
)

var addrs = []string{"localhost:8080", "localhost:8081"}

var allconns = make(map[string]*websocket.Conn)

func DialOne(addr string) {

	conned := false
	for {
		if !conned {
			u := url.URL{Scheme: "ws", Host: addr, Path: "/echo"}
			c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Printf("尝试连接节点%s失败,3s后重试:\n", addr)
				conned = false
				time.Sleep(3 * time.Second)
				continue
			}
			conned = true
			allconns[addr] = c
			fmt.Printf("已创建连接节点%s\n", addr)
		}

	}
}

func DLocal() {
	for _, addr := range addrs {
		go DialOne(addr)
	}
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
	if err != nil {
		fmt.Println("节点消息投送失败")
		panic(err)

	}
	rv := rusult{
		Code: 200,
		Msg:  "消息已成功发送到执行节点",
	}
	err = json.NewEncoder(w).Encode(rv)
	if err != nil {
		//... handle error
		panic(err)
	}

}
