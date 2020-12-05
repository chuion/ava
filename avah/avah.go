package avah

import (
	"ava/core"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{} // use default options

type Task struct {
	Route string `json:"route"`
	Cmd   string `json:"cmd"`
	Args  string `json:"args"`
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	fmt.Print("接到管理端ws接成功\n")
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	p := Task{}
	for {
		err := c.ReadJSON(&p)
		if err != nil {
			log.Println("readjson:", err)
			break
		}
		//log.Printf("recv: %s", message)

		go core.Executor(p.Cmd, p.Args)
		//err = c.WriteMessage(mt, message)
		//if err != nil {
		//	log.Println("write:", err)
		//	break
		//}
	}
}

func HLocal() {
	addr:=strings.Join([]string{"0.0.0.0",":", core.WsPort} , "")
	go Socks5h()
	http.HandleFunc("/echo", echo)
	fmt.Printf("ws监听地址: %s \n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
