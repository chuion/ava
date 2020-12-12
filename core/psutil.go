package core

import (
	"os"
)

func GetPcInfo() (info PcInfo) {
	info.Pid = os.Getpid()
	info.Version = Version
	return info
}
