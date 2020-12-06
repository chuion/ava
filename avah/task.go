package avah

import (
	"ava/core"
	"fmt"
)

func taskrouter(p core.TaskMsg)  {
	//go Executor(p.Cmd, p.Args)
	fmt.Printf("@@@",p)
}
