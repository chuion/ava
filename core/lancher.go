package core

import (
	"log"
	"os/exec"
)


func Executor(name string, arg ...string) {
	log.Printf("启动器接到命令: %s: %s\n", name,arg)
	cmd := exec.Command(name, arg...)
	// 使用CombinedOutput 将stdout stderr合并输出
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("测试执行 failed %s\n", err)
	}
	log.Println("测试执行 标准输出: ", string(out))

	//pids <- cmd.Process.Pid


}
