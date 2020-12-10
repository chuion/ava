package avah

import (
	"bufio"
	"context"
	"github.com/phuslu/log"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func executor(command, arg, taskid, dir string) {
	//go reaper.Reap()

	ctx, cancel := context.WithCancel(context.Background())
	filename := fileConfig(dir, arg)

	script := strings.Split(command, " ")
	log.Debug().Msgf("启动器接到命令: %s %s %s %s\n", script[0], script[1], "placeholder", filename)

	cmd := exec.CommandContext(ctx, script[0], script[1], "placeholder", filename)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Debug().Msgf("程序执行失败 %s", err)
	}
	//fmt.Printf("----------%s 标准输出-----------:\n%s\n", dir, string(out))
	log.Debug().Msgf("程序%s %s成功启动,任务id: %s 进程id: %s", script[0], script[1], taskid, cmd.Process.Pid)
	taskid = taskid + ".log"
	logfile := filepath.Join(dir, taskid)
	writelog(logfile, out)



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

func writelog(filename string, logmsg []byte) {

	var file *os.File
	var err error
	if Exists(filename) {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Debug().Msgf("文件%s打开失败", filename, err)
		}
	} else {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Debug().Msgf("文件%s创建失败", filename, err)
		}
	}

	defer file.Close()
	write := bufio.NewWriter(file)
	write.Write(logmsg)
	write.Flush()

}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
