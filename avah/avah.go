package avah

import (
	"ava/core"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)


var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	fmt.Print("接到管理端ws接成功\n")
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		//log.Printf("recv: %s", message)


		go core.Executor(string(message[:]))
		//err = c.WriteMessage(mt, message)
		//if err != nil {
		//	log.Println("write:", err)
		//	break
		//}
	}
}



func HLocal(addr string) {

	go Socks5h()

	http.HandleFunc("/echo", echo)
	fmt.Printf("ws监听地址: %s \n",addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
