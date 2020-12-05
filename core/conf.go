package core

const WsPort = "4560"
const TcpPort = "4561"
const SocksPort = "4562"
const Web = "4000"


type LauncherConf struct {
	Worker  string `json:"worker"`
	Command string `json:"command"`
}