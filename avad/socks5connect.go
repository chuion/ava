package avad

import (
	socks5 "github.com/armon/go-socks5"
	"github.com/hashicorp/yamux"
	"github.com/phuslu/log"
	"net"
	"strings"
	"time"
)

func connectForSocks(address string) {
	var session *yamux.Session
	server, _ := socks5.New(&socks5.Config{})
	host := strings.Split(address, ":")[0]
	tcpStatus.Set(host, false)

	for {
		status, _ := tcpStatus.Get(host)
		if !status.(bool) {
			for {
				conn, err := net.Dial("tcp", address)
				if err != nil {
					log.Debug().Msgf("连接远端tcp通道 %s失败,%s后重试", address, pongWait)
					time.Sleep(pongWait)
					continue

				}
				tcpStatus.Set(host, true)
				log.Debug().Msgf("已创建连接节点tcp反弹通道%s\n", address)
				session, err = yamux.Server(conn, nil)

				if err != nil {
					//todo 这里的处理好像还有坑
					session.Close()
					panic(err)

				}
				for {
					stream, err := session.Accept()
					if err != nil {
						log.Debug().Msgf("公网节点无法连接%s可能已经关闭", host)
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

		}
	}
}
