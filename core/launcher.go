package core

import (
	"github.com/phuslu/log"
	"os/exec"
)

func Executor(name string, arg ...string) {
	log.Debug().Msgf("启动器接到命令: %s: %s\n", name, arg)
	cmd := exec.Command(name, arg...)
	// 使用CombinedOutput 将stdout stderr合并输出
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Debug().Msgf("测试执行 failed %s\n", err)
	}
	log.Debug().Msgf("测试执行 标准输出: %s", string(out))


}
