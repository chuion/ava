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
	"time"
)

var proxytout = time.Millisecond * 1000 //timeout for wait magicbytes

// listen for agents
func listenForAgents() {
	address := strings.Join([]string{"0.0.0.0", ":", core.TcpPort}, "")

	var err, erry error
	var session *yamux.Session
	var ln net.Listener
	log.Printf("Listening for agents on %s", address)
	ln, err = net.Listen("tcp", address)

	if err != nil {
		log.Printf("Error listening on %s: %v", address, err)

	}

	conn, err := ln.Accept()
	conn.RemoteAddr()
	agentstr := conn.RemoteAddr().String()
	log.Printf("[%s] Got a connection from %v: ", agentstr, conn.RemoteAddr())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errors accepting!")
	}
	conn.SetReadDeadline(time.Now().Add(proxytout))
	log.Printf("[%s] Got Client from %s", agentstr, conn.RemoteAddr())
	conn.SetReadDeadline(time.Now().Add(100 * time.Hour))
	//Add connection to yamux
	session, erry = yamux.Client(conn, nil)
	if erry != nil {
		log.Printf("[%s] Error creating client in yamux for %s: %v", agentstr, conn.RemoteAddr(), erry)

	}

	go listenForClients(agentstr, session)

}

// Catches local clients and connects to yamux
func listenForClients(agentstr string, session *yamux.Session) error {

	address := strings.Join([]string{"0.0.0.0", ":", core.SocksPort}, "")

	log.Printf("[%s] Waiting for clients on %s", agentstr, address)
	ln, err := net.Listen("tcp", address)
	if err!=nil{
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("[%s] Error accepting on %s: %v", agentstr, address, err)
			return err
		}
		if session == nil {
			log.Printf("[%s] Session on %s is nil", agentstr, address)
			conn.Close()
			continue
		}
		log.Printf("[%s] Got client. Opening stream for %s", agentstr, conn.RemoteAddr())

		stream, err := session.Open()
		if err != nil {
			log.Printf("[%s] Error opening stream for %s: %v", agentstr, conn.RemoteAddr(), err)
			return err
		}

		// connect both of conn and stream

		go func() {
			log.Printf("[%s] Starting to copy conn to stream for %s", agentstr, conn.RemoteAddr())
			io.Copy(conn, stream)
			conn.Close()
			log.Printf("[%s] Done copying conn to stream for %s", agentstr, conn.RemoteAddr())
		}()
		go func() {
			log.Printf("[%s] Starting to copy stream to conn for %s", agentstr, conn.RemoteAddr())
			io.Copy(stream, conn)
			stream.Close()
			log.Printf("[%s] Done copying stream to conn for %s", agentstr, conn.RemoteAddr())
		}()
	}
}
