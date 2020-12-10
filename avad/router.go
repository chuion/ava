package avad

import (
	"ava/core"
	"fmt"
	"github.com/phuslu/log"
	"math/rand"
	"time"
)

func fixed(p core.TaskMsg) (rv *result) {
	log.Debug().Msgf("接到 %s的定点任务", p.Route)
	_, err := netAvailable(p.Route)
	if err != nil {
		return &result{400, err.Error()}
	}
	host := workerAvailable(p.Route, p.Worker)
	if host != "" {
		code, msg := send(p.Route, p)
		return &result{code, msg}
	}
	msg := fmt.Sprintf("未在主机: %s上找到业务%s,请检查参数", p.Route, p.Worker)
	return &result{400, msg}
}

func balance(p core.TaskMsg) (rv *result) {
	log.Debug().Msgf("自动解析: %s任务的运行节点", p.Worker)
	if hosts, ok := workerMap[p.Worker]; !ok {
		msg := fmt.Sprintf("投送的业务: %s未找到", p.Worker)
		return &result{400, msg}
	} else {
		host := randone(hosts, p)
		if host == "" {
			msg := fmt.Sprintf("任务: %s在节点 %s都有部署,全节点不可达", p.Route, hosts)
			return &result{400, msg}
		}

		log.Debug().Msgf("任务: %s在节点 %s都有部署,随机投送到: %s执行", p.Route, hosts, host)
		code, msg := send(host, p)
		msg = fmt.Sprintf("任务%s在%s都有部署,随机%s", p.Worker, hosts, msg)
		return &result{code, msg}
	}
}

func randone(hosts []string, p core.TaskMsg) (host string) {
	tmp := make([]string, len(hosts), len(hosts))
	copy(tmp, hosts)
	for range tmp {
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s) // initialize local pseudorandom generator
		index := r.Intn(len(tmp))
		host = tmp[index]
		if _, err := netAvailable(host); err != nil {
			log.Debug().Msgf("任务: %s在节点 %s都有部署,随机节点: %s不可用,更换下一个", p.Worker, hosts, host)
			tmp = append(tmp[:index], tmp[index+1:]...)
			continue
		}
		return host
	}
	return ""

}
