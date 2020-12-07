package avah

import (
	"context"
	"github.com/phuslu/log"
	"os/exec"
	"strings"
)

func executor(command, arg string) {
	//go reaper.Reap()

	ctx, cancel := context.WithCancel(context.Background())

	cmdStr := strings.Join([]string{command, " ", arg}, "")
	script := strings.Split(command, " ")
	log.Debug().Msgf("启动器接到命令: %s\n", cmdStr)
	cmd := exec.CommandContext(ctx, script[0], script[1], arg)
	//cmd.Dir = "/usr/bin"

	err := cmd.Start()
	//out, err := cmd.CombinedOutput()
	if err != nil {
		log.Debug().Msgf("测试执行 failed %s\n", err)
	}
	//log.Debug().Msgf("测试执行 标准输出: %s", out)

	go func() {
		cmd.Wait()
		cancel()
		log.Debug().Msgf("任务: %s执行完成,退出", cmdStr)
	}()

}
