package core

import (
	cmap "github.com/orcaman/concurrent-map"
	"time"
)

const WsPort = "4560"
const TcpPort = "4561"
const SocksPort = "4562"
const Web = "4000"

type LauncherConf struct {
	Worker  string `json:"worker"`
	Command string `json:"command"`
	Dir     string `json:"dir"`
	PcInfo
}

var ProcessStatus = cmap.New()

type ProcessInfo struct {
	TaskId  string
	Pid     int32
	Mem     uint64
	Threads int32
	Files   int
	CpuPer  float64
}

type PcInfo struct {
	Version   string `json:"version"`
	ProNum   int
	MemTotal uint64
	MemUsed  uint64
	TotalPercent float64
	ProStatus []ProcessInfo
}

type TaskMsg struct {
	Route  string `json:"route"`
	Worker string `json:"worker"`
	TaskID string `json:"task_id"`
	Params string `json:"params"`
}

const PongWait = 20 * time.Second
const UpdateWait = 5 * time.Second

//白名单
var Sites []string

var Version = "1.1"
