package avah

import (
	"ava/core"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/phuslu/log"
	"github.com/spf13/viper"
)

func initWorker() (config core.LauncherConf) {

	viper.SetConfigName("launcher") // 设置配置文件名 (不带后缀)
	viper.AddConfigPath(".")        // 第一个搜索路径
	err := viper.ReadInConfig()     // 读取配置数据
	if err != nil {
		panic(fmt.Errorf("未找到launcher.json: %s \n", err))
	}
	viper.Unmarshal(&config) // 将配置信息绑定到结构体上
	return
}

func infoReg(c *websocket.Conn) {
	log.Debug().Msgf("接到管理端ws连接成功,开始注册/更新业务功能\n")
	config := initWorker()
	err := c.WriteJSON(config)
	if err != nil {
		log.Debug().Msgf("注册/更新业务功能失败\n")

	}
}
