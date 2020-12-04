package main

import (
	"ava/avad"
	"ava/avah"
	"flag"
	"fmt"
)


func main() {
	var flags struct {
		Mode string
		Addr string
	}

	flag.StringVar(&flags.Mode, "m", "h", "程序哪种模式运行d/h(打印版本)")
	flag.StringVar(&flags.Addr, "addr", "localhost:8080", "以h运行时,监听的地址和端口,默认 localhost:8080")
	flag.Parse()

	if flags.Mode == "d" {
		fmt.Println("程序启动以管理模式运行")
		avad.DLocal()
	}

	if flags.Mode == "h" {
		fmt.Println("程序启动以节点模式运行")
		avah.HLocal(flags.Addr)

	}



}
