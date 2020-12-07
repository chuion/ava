package avad

import (
	"ava/core"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/phuslu/log"
	"net/http"
)

type rusult struct {
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
		code, msg := send(p.Route, p)
		json.NewEncoder(w).Encode(&rusult{code, msg})
		return
	}

	log.Debug().Msgf("自动解析: %s任务的运行节点", p.Worker)
	if host, ok := workerMap[p.Worker]; !ok {
		msg = fmt.Sprintf("投送的业务: %s未找到", p.Worker)
		json.NewEncoder(w).Encode(&rusult{400, msg})
		return
	} else {
		code, msg := send(host, p)
		json.NewEncoder(w).Encode(&rusult{code, msg})
		return
	}
}

func nodeAvailable(host string) (conn *websocket.Conn, err error) {
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
	conn, err := nodeAvailable(host)
	if conn != nil {
		p.Command = workerCommand[p.Worker] //补充上command

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
