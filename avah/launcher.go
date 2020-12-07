package avah

import (
	"context"
	"fmt"
	"github.com/phuslu/log"
	"os/exec"
	"strings"
)

func executor(command, arg, dir string) {
	//go reaper.Reap()

	ctx, cancel := context.WithCancel(context.Background())

	script := strings.Split(command, " ")
	log.Debug().Msgf("启动器接到命令: %s %s\n", command, arg)
	cmd := exec.CommandContext(ctx, script[0], script[1], arg)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Debug().Msgf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	go func() {
		cmd.Wait()
		cancel()
		log.Debug().Msgf("任务: %s执行完成,退出", command)
	}()

}
