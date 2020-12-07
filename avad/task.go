package avad

import (
	"ava/core"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/phuslu/log"
	"math/rand"
	"net/http"
	"time"
)

type result struct {
	Code int
	Msg  string
}

func taskRouter(w http.ResponseWriter, r *http.Request) {
	var p core.TaskMsg
	var msg string
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if p.Route != "" {
		log.Debug().Msgf("接到 %s的定点任务", p.Route)
		host := workerAvailable(p.Route, p.Worker)
		if host != "" {
			code, msg := send(p.Route, p)
			json.NewEncoder(w).Encode(&result{code, msg})
			return
		} else {
			msg = fmt.Sprintf("未在主机: %s上找到业务,或者节点不存在: %s,请检查配置", p.Route, p.Worker)
			json.NewEncoder(w).Encode(&result{400, msg})
			return
		}
	}

	log.Debug().Msgf("自动解析: %s任务的运行节点", p.Worker)
	if hosts, ok := workerMap[p.Worker]; !ok {
		msg = fmt.Sprintf("投送的业务: %s未找到", p.Worker)
		json.NewEncoder(w).Encode(&result{400, msg})
		return
	} else {
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s) // initialize local pseudorandom generator
		t := r.Intn(len(hosts))
		host := hosts[t]
		log.Debug().Msgf("任务: %s在节点 %s都有部署,随机投送到: %s执行", p.Route, hosts, host)
		code, msg := send(host, p)
		msg = fmt.Sprintf("任务%s在%s都有部署,随机%s",p.Worker,hosts, msg)
		json.NewEncoder(w).Encode(&result{code, msg})
		return
	}
}

func workerAvailable(host, workerdst string) (hostdst string) {
	if workers, ok := workerMapR[host]; ok {
		for _, worker := range workers {
			if worker == workerdst {
				return host
			}
		}
		return ""
	}
	return ""
}

func netAvailable(host string) (conn *websocket.Conn, err error) {
	if conn, ok := wsConns[host]; !ok {
		return nil, fmt.Errorf("未找到节点: %s,请检查输入", host)
	} else {
		status, _ := wsStatus.Get(host)
		if !status.(bool) {
			return nil, fmt.Errorf("节点: %s,网络中断", host)
		} else {
			return conn, nil
		}
	}

}

func send(host string, p core.TaskMsg) (code int, msg string) {
	conn, err := netAvailable(host)
	if conn != nil {
		log.Debug().Msgf("发送前原始参数: %s", p)
		err = conn.WriteJSON(p)
		if err != nil {
			log.Debug().Msgf("投送失败,节点: %s可能已不在线", host)
			return 400, "投送失败,节点可能已不在线"
		}
		return 200, fmt.Sprintf("投送到: %s成功", host)

	} else {
		return 200, fmt.Sprintf("%s", err)
	}

}
