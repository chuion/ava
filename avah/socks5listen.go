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
var socksListen net.Listener
var lis = false
var session *yamux.Session

// listen for agents
func listenForAgents() {
	address := strings.Join([]string{"0.0.0.0", ":", core.TcpPort}, "")

	var err, erry error

	var ln net.Listener
	log.Printf("Listening for agents on %s", address)
	ln, err = net.Listen("tcp", address)

	if err != nil {
		log.Printf("Error listening on %s: %v", address, err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Errors accepting!")
			return
		}
		//defer conn.Close()

		agentStr := conn.RemoteAddr().String()
		log.Printf("[%s] Got a connection from %v: ", agentStr, conn.RemoteAddr())
		//_ = conn.SetReadDeadline(time.Now().Add(proxytout))
		//conn.SetReadDeadline(time.Now().Add(100 * time.Hour))

		//Add connection to yamux
		session, erry = yamux.Client(conn, nil)
		if erry != nil {
			log.Printf("[%s] Error creating client in yamux for %s: %v", agentStr, conn.RemoteAddr(), erry)
		}

		go listenForClients(agentStr)

	}

}

// Catches local clients and connects to yamux


func listenForClients(agentStr string) error {
	var err error
	address := strings.Join([]string{"0.0.0.0", ":", core.SocksPort}, "")

	if !lis {
		log.Printf("[%s] Waiting for clients on %s", agentStr, address)
		socksListen, err = net.Listen("tcp", address)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Errors accepting!")
		}
		lis = true
	}

	for {
		conn, err := socksListen.Accept()
		log.Printf("接到代理请求")
		if err != nil {
			log.Printf("[%s] Error accepting on %s: %v", agentStr, address, err)
			return err
		}
		if session == nil {
			log.Printf("[%s] Session on %s is nil", agentStr, address)
			conn.Close()
			continue
		}
		log.Printf("[%s] Got client. Opening stream for %s", agentStr, conn.RemoteAddr())

		stream, err := session.Open()
		if err != nil {
			log.Printf("[%s] Error opening stream for %s: %v", agentStr, conn.RemoteAddr(), err)
			_ = session.Close()

			return err
		}

		go relay(conn, stream)

	}
}

func relay(conn, stream net.Conn) {
	type res struct {
		N   int64
		Err error
	}
	ch := make(chan res)

	go func() {
		n, err := io.Copy(stream, conn)
		_ = stream.SetDeadline(time.Now()) // wake up the other goroutine blocking on stream
		_ = conn.SetDeadline(time.Now())   // wake up the other goroutine blocking on conn
		defer stream.Close()
		defer conn.Close()
		ch <- res{n, err}
	}()

	_, err := io.Copy(conn, stream)
	defer stream.Close()
	defer conn.Close()
	_ = stream.SetDeadline(time.Now()) // wake up the other goroutine blocking on stream
	_ = conn.SetDeadline(time.Now())   // wake up the other goroutine blocking on left
	rs := <-ch

	if err == nil {
		err = rs.Err
	}

}
