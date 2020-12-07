package core

const WsPort = "4560"
const TcpPort = "4561"
const SocksPort = "4562"
const Web = "4000"

type LauncherConf struct {
	Worker  string `json:"worker"`
	Command string `json:"command"`
}

type TaskMsg struct {
	Route   string `json:"route"`
	Worker  string `json:"worker"`
	TaskID  string `json:"task_id"`
	Params  string `json:"params"`
	Command string `json:"command"`
}
