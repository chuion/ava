package core

import (
	"github.com/shirou/gopsutil/v3/process"
)

func GetPcInfo() (info PcInfo) {
	ch := ProcessStatus.IterBuffered()

	var rv []ProcessInfo
	for item := range ch {
		pid, ok := item.Val.(int)
		if !ok {
			continue
		}
		proIns := process.Process{Pid: int32(pid)}
		Threads, _ := proIns.NumThreads()
		t, _ := proIns.MemoryInfo()
		f, _ := proIns.OpenFiles()

		rv = append(rv, ProcessInfo{
			TaskId:  item.Key,
			Mem:     t.RSS,
			Pid:     int32(pid),
			Threads: Threads,
			Files:   len(f),
		})
	}
	info.ProStatus = rv
	info.Version = Version
	return info
}
