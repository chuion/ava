package avad

import (
	"ava/core"
	"github.com/gorilla/websocket"
	"github.com/phuslu/log"
	"strings"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.

	// Send pings to peer with this period. Must be less than PongWait.
	pingPeriod = (core.PongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

func ping() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for host, ws := range wsConns {
				if ws == nil {
					reconnect(host)
					continue
				}

				ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(core.PongWait)); return nil })
				err := ws.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					log.Debug().Msgf("节点 %s的ws心跳检测失败,触发重新连接:%s", host, err)
					reconnect(host)

				} else {
					log.Debug().Msgf("节点 %s的ws心跳检测正常", host)
				}
			}
		}
	}
}

func reconnect(host string) {
	addrWs := strings.Join([]string{host, ":", core.WsPort}, "")
	addrTcp := strings.Join([]string{host, ":", core.TcpPort}, "")
	go dialWs(addrWs)
	go connectForSocks(addrTcp)
}
