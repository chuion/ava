package core

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"strconv"
	"time"
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

		cpuper, _ := proIns.CPUPercent()
		cpuper1, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", cpuper), 64)

		Threads, _ := proIns.NumThreads()
		t, _ := proIns.MemoryInfo()
		f, _ := proIns.OpenFiles()

		rv = append(rv, ProcessInfo{
			TaskId:  item.Key,
			Mem:     t.RSS,
			Pid:     int32(pid),
			Threads: Threads,
			Files:   len(f),
			CpuPer:  cpuper1,
		})
	}
	info.ProStatus = rv
	info.Version = Version
	info.ProNum = len(rv)
	m, _ := mem.VirtualMemory()
	info.MemTotal = m.Total
	info.MemUsed = m.Used
	cpu, _ := cpu.Percent(3*time.Second, false)
	cpu1, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", cpu[0]), 64)
	info.TotalPercent = cpu1

	return info
}
