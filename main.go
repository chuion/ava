package main

import (
	"ava/avad"
	"ava/avah"
	"fmt"
	"github.com/spf13/viper"
	"os"
)



func main() {

	if len(os.Args) > 1 {
		fmt.Printf("程序启动以管理模式运行,配置文件为: %s\n", os.Args[1])
		nodes := LoadConfig(os.Args[1])
		avad.DLocal(nodes)
	}

	fmt.Println("程序启动以节点模式运行")
	avah.HLocal()

}

func LoadConfig(config string) []string {

	viper.SetConfigFile(config)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("配置文件读取失败: %s\n", err)
		os.Exit(1)
	}
	nodes := viper.GetStringSlice("nodes")

	return nodes
}
