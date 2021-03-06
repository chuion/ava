package main

import (
	"ava/avad"
	"ava/avah"
	"ava/core"
	"github.com/phuslu/log"
	"github.com/spf13/viper"
	baseLog "log"
	"net/http"
	"os"
	"runtime"

	_ "net/http/pprof"
)

func init() {

	if log.IsTerminal(os.Stderr.Fd()) {
		log.DefaultLogger = log.Logger{
			Caller: 1,
			Writer: &log.ConsoleWriter{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: true,
			},
		}
	}

}

func main() {
	runtime.GOMAXPROCS(1)
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)
	go func() {
		baseLog.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	if len(os.Args) > 1 {
		log.Debug().Msgf("程序启动以管理模式运行,配置文件为: %s\n", os.Args[1])
		nodes := LoadConfig(os.Args[1])
		avad.Manger(nodes)
	}

	log.Debug().Msgf("程序启动以节点模式运行")
	avah.Node()

}

func LoadConfig(config string) []string {

	viper.SetConfigFile(config)
	err := viper.ReadInConfig()
	if err != nil {
		log.Debug().Msgf("配置文件读取失败: %s\n", err)
		os.Exit(1)
	}
	nodes := viper.GetStringSlice("nodes")
	core.Sites = viper.GetStringSlice("sites")
	core.PerMachineProcess = viper.GetInt("permachineprocess")

	return nodes
}
