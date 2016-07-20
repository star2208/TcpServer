package tcpserver

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

// TcpServer 类，虚拟一个Socket服务器
type TcpServer struct {
	// 取名
	Name string
	// 监听端口
	ListenPort int16
	// 停止监听
	stop_listener chan bool
	// 监听已经停止
	listener_is_stopped chan bool
}

func NewTcpServer(name string, listen_port int16) *TcpServer {
	fmt.Println("Create TcpServer name:", name, "port: ", listen_port)
	return &TcpServer{name, listen_port, nil, nil}
}

func (s *TcpServer) handleConnection(conn net.Conn) {
	io.Copy(conn, conn)
	// Shut down the connection.
	conn.Close()
}

func (s *TcpServer) Listen(stop_listener, listener_is_stopped chan bool) (err error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+strconv.Itoa(int(s.ListenPort)))
	if err != nil {
		fmt.Println("Started Listener fail.", err.Error())
	}

	tcplistener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Started Listener fail.", err.Error())
	}
	fmt.Println("Listener Started...", "Port:", tcpAddr.Port)

	isStop := false
	for {

		select {
		case <-s.stop_listener:
			isStop = true
		default:
		}

		if isStop {
			break
		}

		tcplistener.SetDeadline(time.Now().Add(1e9))

		conn, err := tcplistener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				//fmt.Println(err)
				continue
			}
			fmt.Println(err)
		} else {
			go s.handleConnection(conn)
		}
	}

	fmt.Println("Listener Stopped.")
	s.listener_is_stopped <- true

	return
}

func (s *TcpServer) Start() (err error) {
	fmt.Println("TcpServer", s.Name, "started.")

	s.stop_listener = make(chan bool)
	s.listener_is_stopped = make(chan bool)

	go s.Listen(s.stop_listener, s.listener_is_stopped)

	return nil
}

func (s *TcpServer) Stop() (err error) {
	s.stop_listener <- true
	<-s.listener_is_stopped
	fmt.Println("TcpServer", s.Name, "stoped.")
	return nil
}
