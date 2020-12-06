package avad

import (
	socks5 "github.com/armon/go-socks5"
	"github.com/hashicorp/yamux"
	"github.com/phuslu/log"
	"net"
)


func connectForSocks(address string) error {
	var session *yamux.Session
	server, err := socks5.New(&socks5.Config{})
	if err != nil {
		return err
	}

	var conn net.Conn

	log.Debug().Msgf("Connecting to far end")

	conn, err = net.Dial("tcp", address)

	if err != nil {
		panic(err)
	}

	log.Debug().Msgf("Starting client")

	session, err = yamux.Server(conn, nil)

	if err != nil {
		return err
	}

	for {
		stream, err := session.Accept()
		log.Debug().Msgf("Accepting stream")
		if err != nil {
			return err
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
