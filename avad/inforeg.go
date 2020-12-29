package avad

import (
	"ava/core"
	"github.com/gorilla/websocket"
	"github.com/phuslu/log"
)

var workerMap = make(map[string][]string)
var workerMapR = make(map[string][]string)
var Ver = make(map[string]core.PcInfo)

func getNodeInfo(host string, c *websocket.Conn) {
	for {
		p := make(map[string]core.LauncherConf)
		err := c.ReadJSON(&p)
		if err != nil {
			log.Debug().Msgf("读取节点: %s注册信息失败 %s", host, err)
			return
		}

		if value, ok := p["info"]; ok {
			Ver[host] = value.PcInfo
			//log.Debug().Msgf("读取节点: %s 状态信息成功", host)
			continue
		}

		for k, _ := range p {
			workerMapR[host] = append(workerMapR[host], k)
			workerMap[k] = append(workerMap[k], host)
		}
		//去重配置
		for k, v := range workerMap {
			workerMap[k] = RemoveRepeatedElement(v)
		}
		for k, v := range workerMapR {
			workerMapR[k] = RemoveRepeatedElement(v)
		}
		log.Debug().Msgf("读取节点: %s注册信息成功", host)
	}

}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}
