package avah

import (
	"ava/core"
	"github.com/gorilla/websocket"
	"github.com/phuslu/log"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{} // use default options


func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Debug().Msgf("ws握手失败: %s", err)
		return
	}
	reg(c)

	defer c.Close()
	p := Task{}
	for {
		err := c.ReadJSON(&p)
		if err != nil {
			log.Debug().Msgf("读取数据失败:", err)
			break
		}
		//log.Printf("recv: %s", message)

		go Executor(p.Cmd, p.Args)
		//err = c.WriteMessage(mt, message)
		//if err != nil {
		//	log.Println("write:", err)
		//	break
		//}
	}
}

func HLocal() {

	go listenForAgents()

	addr := strings.Join([]string{"0.0.0.0", ":", core.WsPort}, "")
	http.HandleFunc("/echo", echo)
	log.Debug().Msgf("ws监听地址: %s \n", addr)
	http.ListenAndServe(addr, nil)
}
