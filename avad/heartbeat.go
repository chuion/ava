package avad

import (
	"ava/core"
	"github.com/gorilla/websocket"
	"github.com/phuslu/log"
	"strings"
	"time"
)

const (
	// Send pings to peer with this period. Must be less than PongWait.
	pingPeriod = (core.PongWait * 9) / 10
)

func ping() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		<-ticker.C
		ch := wsConns.IterBuffered()
		for item := range ch {
			host := item.Key
			ws := item.Val.(*websocket.Conn)

			if ws == nil {
				reconnect(host)
				continue
			}

			ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(core.PongWait)); return nil })
			err := ws.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Debug().Msgf("节点 %s的ws心跳检测失败,触发重新连接:%s", host, err)
				reconnect(host)
			}
			log.Debug().Msgf("节点 %s的ws心跳检测正常", host)
		}
	}
}

func reconnect(host string) {
	addrWs := strings.Join([]string{host, ":", core.WsPort}, "")
	addrTcp := strings.Join([]string{host, ":", core.TcpPort}, "")
	go dialWs(addrWs)
	go dialTcp(addrTcp)
}
