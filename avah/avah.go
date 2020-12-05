package avah

import (
	"ava/core"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"

	"net/http"
	"strings"
	"github.com/phuslu/log"
)

var upgrader = websocket.Upgrader{} // use default options

type Task struct {
	Route string `json:"route"`
	Cmd   string `json:"cmd"`
	Args  string `json:"args"`
}

type Launcher struct {


}

func initWorker() (config Launcher) {

	viper.SetConfigName("launcher") // 设置配置文件名 (不带后缀)
	viper.AddConfigPath(".")        // 第一个搜索路径
	err := viper.ReadInConfig()     // 读取配置数据
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.Unmarshal(&config) // 将配置信息绑定到结构体上
	fmt.Println(config)
	return
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	log.Debug().Msgf("接到管理端ws接成功\n")
	if err != nil {
		log.Debug().Msgf("upgrade:", err)
		return
	}
	defer c.Close()
	p := Task{}
	for {
		err := c.ReadJSON(&p)
		if err != nil {
			log.Debug().Msgf("读取数据失败:", err)
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
	addr := strings.Join([]string{"0.0.0.0", ":", core.WsPort}, "")
	go Socks5h()
	http.HandleFunc("/echo", echo)
	log.Debug().Msgf("ws监听地址: %s \n", addr)
	http.ListenAndServe(addr, nil)
}
