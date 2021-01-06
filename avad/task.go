package avad

import (
	"ava/core"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/phuslu/log"
	"net/http"
)

type result struct {
	Code int
	Msg  string
	Route string
}

func taskRouter(w http.ResponseWriter, r *http.Request) {
	var p core.TaskMsg
	var rv = &result{}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if p.Route != "" {
		rv = fixed(p)

	} else {
		rv = balance(p)
	}
	json.NewEncoder(w).Encode(rv)

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

	ws, ok := wsConns.Get(host)
	if !ok {
		return nil, fmt.Errorf("未找到节点: %s,请检查输入", host)
	}
	conn = ws.(*websocket.Conn)
	status, _ := wsStatus.Get(host)
	if !status.(bool) {
		return nil, fmt.Errorf("节点: %s,网络中断", host)
	}

	return conn, nil

}

func send(host string, p core.TaskMsg) (code int, msg string) {
	conn, err := netAvailable(host)
	if err != nil {
		return 400, fmt.Sprintf("%s", err)
	}

	if conn != nil {
		log.Debug().Msgf("发送前原始参数: %s  %s  %s", p.Worker, p.Route, p.TaskID)
		err = conn.WriteJSON(p)
		if err != nil {
			log.Debug().Msgf("投送失败,节点: %s可能已不在线", host)
			wsStatus.Set(host, false)
			return 400, fmt.Sprintf("投送失败,节点: %s可能已不在线 %s", host, err)
		}
	}
	return 200, fmt.Sprintf("投送到: %s成功", host)
}
