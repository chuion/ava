package avah

import (
	"ava/core"
	"github.com/phuslu/log"
)

func taskrouter(p core.TaskMsg) {
	cmd := allConfig[p.Worker].Command
	log.Debug().Msgf("接收到的原始参数: %s", p)
	dir := allConfig[p.Worker].Dir
	go executor(cmd, p.Params, dir)

}
