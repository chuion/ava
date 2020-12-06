package avah

import "ava/core"

func taskrouter(p core.Task)  {
	go Executor(p.Cmd, p.Args)
}
