package lancher

import (
	"log"
)


func Executor(name string, arg ...string) {
	//cmd := exec.Command(name, arg...)
	//err := cmd.Start()
	//if err != nil {
	//	fmt.Printf("Error starting  command: %s\n", arg)
	//	return
	//}
	//
	//pids <- cmd.Process.Pid

	log.Printf("启动器接到命令: %s \n", name)

}
