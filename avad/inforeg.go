package avad

import (
	"ava/core"
	"github.com/gorilla/websocket"
	"github.com/phuslu/log"
)

var workerMap = make(map[string]string)
var workerMapR = make(map[string]string)
var workerCommand = make(map[string]string)

func infoReg(host string, c *websocket.Conn) {
	p := core.LauncherConf{}
	err := c.ReadJSON(&p)
	if err != nil {
		log.Debug().Msgf("读取节点: %s注册信息失败", host)
	}
	log.Debug().Msgf("接收节点: %s注册信息成功,可运行%s", host, p.Worker)
	workerMap[p.Worker] = host
	workerMapR[host] = p.Worker
	workerCommand[p.Worker] = p.Command

}
