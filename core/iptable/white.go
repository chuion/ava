package iptable

import (
	"github.com/phuslu/log"
	"github.com/spf13/viper"
	"os"
)

var allList []string

func init() {
	if len(os.Args) != 2 {
		return
	}
	viper.SetConfigFile("allow.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Debug().Msgf("白名单配置文件读取失败: %s\n", err)
		return
	}
	allList = viper.GetStringSlice("sites")
	log.Debug().Msgf("白名单配置文件读取成功")
}

func Allow(dst string) bool {
	if !stringInSlice(dst, allList) {
		return false
	}
	return true
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
