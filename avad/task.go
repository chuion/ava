package avad

import (
	"encoding/json"
	"fmt"
	"github.com/phuslu/log"
	"net/http"
)

type Task struct {
	Route string
	Cmd   string
	Args  string
}

type rusult struct {
	Code int
	Msg  string
}

func taskrouter(w http.ResponseWriter, r *http.Request) {
	var p Task
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//todo 往哪个连接发应解析业务逻辑
	addr := p.Route
	msg := "投送成功"
	if c, ok := allconn[addr]; ok {
		err = c.WriteJSON(p)
		if err != nil {
			log.Debug().Msgf("投送失败,节点可能已不在线")
			msg = "投送失败,节点可能已不在线"
		}
	} else {
		msg = fmt.Sprintf("未找到%s对应的socket连接", p.Route)
	}

	rv := rusult{
		Code: 200,
		Msg:  msg,
	}
	err = json.NewEncoder(w).Encode(rv)
	if err != nil {
		//... handle error
		panic(err)
	}

}
