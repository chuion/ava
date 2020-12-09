package avad

import (
	"ava/core"
	"ava/core/go-socks5"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/yamux"
	"github.com/phuslu/log"
	"net"
	"net/url"
	"strings"
)

func dialTcp(address string) {
	var session *yamux.Session
	server, _ := socks5.New(&socks5.Config{})
	host := strings.Split(address, ":")[0]
	tcpStatus.Set(host, false)

	status, _ := tcpStatus.Get(host)
	if !status.(bool) {
		for {
			conn, err := net.Dial("tcp", address)
			if err != nil {
				log.Debug().Msgf("连接远端tcp通道 %s失败,%s后重试", address, core.PongWait)
				return
			}
			tcpStatus.Set(host, true)
			log.Debug().Msgf("已创建连接节点tcp反弹通道%s\n", address)
			session, err = yamux.Server(conn, nil)

			if err != nil {
				//todo 这里的处理好像还有坑
				session.Close()
				panic(err)

			}
			relay(host, session, server)
		}

	}
}

func dialWs(addr string) {
	host := strings.Split(addr, ":")[0]
	wsStatus.Set(host, false)
	status, _ := wsStatus.Get(host)
	if !status.(bool) {
		u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Debug().Msgf("尝试连接节点ws通道%s失败,%s后重试:\n", addr, core.PongWait)
			return
		}
		host := strings.Split(u.Host, ":")[0]
		wsStatus.Set(host, true)
		//wsConns[host] = c
		wsConns.Set(host, c)
		log.Debug().Msgf("已创建连接节点ws通道%s\n", addr)
		//读取注册信息
		infoReg(host, c)
	}
}

func relay(host string, session *yamux.Session, server *socks5.Server) {
	for {
		stream, err := session.Accept()
		if err != nil {
			log.Debug().Msgf("公网节点无法连接%s可能已经关闭", host)
			wsStatus.Set(host, false)
			tcpStatus.Set(host, false)
			break
		}
		log.Debug().Msgf("Passing off to socks5")
		go func() {
			err = server.ServeConn(stream)
			if err != nil {
				log.Debug().Err(err)

			}
		}()
	}

}
