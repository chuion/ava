package avah

import (
	"ava/core"
	"fmt"
	"github.com/hashicorp/yamux"
	"github.com/phuslu/log"
	"io"
	"net"
	"os"
	"strings"
)

var session *yamux.Session

// Catches yamux connecting to us
func listenForTcp() {
	address := strings.Join([]string{"0.0.0.0", ":", core.TcpPort}, "")

	log.Debug().Msgf("Listening for the far end")
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return
	}
	for {
		conn, err := ln.Accept()
		log.Debug().Msgf("Got a client")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Errors accepting!")
		}
		// Add connection to yamux
		session, err = yamux.Client(conn, nil)
	}
}

// Catches clients and connects to yamux
func listenForClients() error {
	address := strings.Join([]string{"0.0.0.0", ":", core.SocksPort}, "")
	log.Debug().Msgf("Waiting for clients on: %s", address)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		// TODO dial socks5 through yamux and connect to conn

		if session == nil {
			conn.Close()
			continue
		}
		log.Debug().Msgf("Got a client")

		log.Debug().Msgf("Opening a stream")
		stream, err := session.Open()
		if err != nil {
			return err
		}

		// connect both of conn and stream

		go func() {
			log.Debug().Msgf("Starting to copy conn to stream")
			io.Copy(conn, stream)
			conn.Close()
		}()
		go func() {
			log.Debug().Msgf("Starting to copy stream to conn")
			io.Copy(stream, conn)
			stream.Close()
			log.Debug().Msgf("Done copying stream to conn")
		}()
	}
}
