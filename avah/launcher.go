package avah

import (
	"context"
	"fmt"
	"github.com/phuslu/log"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func executor(command, arg, dir string) {
	//go reaper.Reap()

	ctx, cancel := context.WithCancel(context.Background())
	filename := fileConfig(dir, arg)

	script := strings.Split(command, " ")
	log.Debug().Msgf("启动器接到命令: %s %s %s %s\n", script[0], script[1],"placeholder",filename)

	cmd := exec.CommandContext(ctx, script[0], script[1],"placeholder",filename)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Debug().Msgf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	go func() {
		cmd.Wait()
		cancel()
		os.Remove(filename)
		log.Debug().Msgf("任务: %s执行完成,退出", command)
	}()

}

func fileConfig(dir, arg string) (filename string) {
	tmpFile, err := ioutil.TempFile(dir, "arg-")
	if err != nil {
		log.Debug().Msgf("创建文件型参数失败: %s", err)
		return ""
	}
	_, err = tmpFile.Write([]byte(arg))
	if err != nil {
		log.Debug().Msgf("向参数文件写入内容失败 %s", err)
		return ""
	}
	return tmpFile.Name()
}
