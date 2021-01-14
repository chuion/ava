package avah

import (
	"ava/core"
	"github.com/gorilla/websocket"
	"github.com/phuslu/log"
	"net/http"
	"strings"
	"time"
)

var upgrader = websocket.Upgrader{} // use default options
var conn *websocket.Conn

func update() {
	ticker := time.NewTicker(core.UpdateWait)
	defer ticker.Stop()
	tmp := make(map[string]core.LauncherConf)
	for {
		tmp["info"] = core.LauncherConf{
			PcInfo: core.GetPcInfo(),
		}
		taskchan <- tmp
		//todo 使用文件监控实现配置文件变更才更新
		listAll(".")
		taskchan <- allConfig
		<-ticker.C
	}
}

func updateInfo() {
	go sendMsg()
	go update()
}

func dial(w http.ResponseWriter, r *http.Request) {
	var err error
	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Debug().Msgf("ws握手失败: %s", err)
		return
	}
	log.Debug().Msgf("接到管理端ws连接")

	updateInfo()
	defer conn.Close()

	p := core.TaskMsg{}
	for {
		err := conn.ReadJSON(&p)
		if err != nil {
			log.Debug().Msgf("读取数据失败,管理节点可能已关闭")
			break
		}
		//接收信息,给到路由分发
		taskrouter(p)
	}
}

func Node() {

	go listenTcp()

	addr := strings.Join([]string{"0.0.0.0", ":", core.WsPort}, "")
	http.HandleFunc("/ws", dial)
	log.Debug().Msgf("ws监听地址: %s ", addr)
	http.ListenAndServe(addr, nil)
}
