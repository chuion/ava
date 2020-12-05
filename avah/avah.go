package avah

import (
	"ava/core"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"

	"github.com/phuslu/log"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{} // use default options
var registed bool = false

type Task struct {
	Route string `json:"route"`
	Cmd   string `json:"cmd"`
	Args  string `json:"args"`
}


func initWorker() (config core.LauncherConf) {

	viper.SetConfigName("launcher") // 设置配置文件名 (不带后缀)
	viper.AddConfigPath(".")        // 第一个搜索路径
	err := viper.ReadInConfig()     // 读取配置数据
	if err != nil {
		panic(fmt.Errorf("未找到launcher.json: %s \n", err))
	}
	viper.Unmarshal(&config) // 将配置信息绑定到结构体上
	//fmt.Println(config)
	return
}

func reg(c *websocket.Conn)  {
		log.Debug().Msgf("接到管理端ws连接成功,开始注册/更新业务功能\n")
		config := initWorker()
		err := c.WriteJSON(config)
		if err != nil {
			log.Debug().Msgf("注册/更新业务功能失败\n")

		}
	}




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

		go core.Executor(p.Cmd, p.Args)
		//err = c.WriteMessage(mt, message)
		//if err != nil {
		//	log.Println("write:", err)
		//	break
		//}
	}
}

func HLocal() {
	go listenForTcp()
	go listenForClients()

	addr := strings.Join([]string{"0.0.0.0", ":", core.WsPort}, "")
	http.HandleFunc("/echo", echo)
	log.Debug().Msgf("ws监听地址: %s \n", addr)
	http.ListenAndServe(addr, nil)
}
