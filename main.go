package main

import (
	"ava/avad"
	"ava/avah"
	"github.com/phuslu/log"
	"github.com/spf13/viper"
	"net/http"
	"os"
	baselog "log"
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
		baselog.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	if len(os.Args) > 1 {
		log.Debug().Msgf("程序启动以管理模式运行,配置文件为: %s\n", os.Args[1])
		nodes := LoadConfig(os.Args[1])
		avad.DLocal(nodes)
	}

	log.Debug().Msgf("程序启动以节点模式运行")
	avah.HLocal()

}

func LoadConfig(config string) []string {

	viper.SetConfigFile(config)
	err := viper.ReadInConfig()
	if err != nil {
		log.Debug().Msgf("配置文件读取失败: %s\n", err)
		os.Exit(1)
	}
	nodes := viper.GetStringSlice("nodes")

	return nodes
}
