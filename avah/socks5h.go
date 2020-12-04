package avah

import (
	"fmt"
	"io"
	"net"
	"time"
)

var ip = "0.0.0.0"
var port = 18080
var s5port = 28080

func Socks5h() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), port, ""})
	if err != nil {
		panic(err)
	}
	defer listen.Close()
	fmt.Printf("tcp监听地址为%s:%d:\n", net.ParseIP(ip),port)
	socks5listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), s5port, ""})
	if err != nil {
		panic(err)
	}
	defer socks5listen.Close()
	fmt.Printf("本机socks5端口为 %d:\n", s5port)

	Server(listen, socks5listen)
}

func Server(listen *net.TCPListener, s5listen *net.TCPListener) {
	for {
		s5conn, err := s5listen.Accept()
		if err != nil {
			fmt.Println("连接socks5服务异常:", err.Error())
			continue
		}
		//fmt.Println("接受客户端连接来自:", s5conn.RemoteAddr().String())

		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("接受被控节点连接异常:", err.Error())
			continue
		}
		//fmt.Println("被控节点连接来自:", conn.RemoteAddr().String())

		go func() {
			defer s5conn.Close()
			defer conn.Close()
			relay(conn, s5conn)
		}()
	}
}

func relay(left, right net.Conn) (int64, int64, error) {
	type res struct {
		N   int64
		Err error
	}
	ch := make(chan res)

	go func() {
		n, err := io.Copy(right, left)
		right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
		left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
		ch <- res{n, err}
	}()

	n, err := io.Copy(left, right)
	right.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
	left.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
	rs := <-ch

	if err == nil {
		err = rs.Err
	}
	return n, rs.N, err
}


