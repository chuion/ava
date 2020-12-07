package avah

import (
	"ava/core"
	"github.com/phuslu/log"
)

func taskrouter(p core.TaskMsg) {
	log.Debug().Msgf("接收到的原始参数: %s", p)
	go executor(p.Command, p.Params)


}
