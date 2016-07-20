package tcpserver

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

// TcpServer 类，虚拟一个Socket服务器
type TcpServer struct {
	// 取名
	Name string
	// 监听端口
	ListenPort int16
	// 停止服务
	stop_server chan bool
	// 服务已停止
	waitGroup *sync.WaitGroup
}

func NewTcpServer(name string, listen_port int16) *TcpServer {
	fmt.Println("Create TcpServer name:", name, "port: ", listen_port)
	return &TcpServer{
		Name:        name,
		ListenPort:  listen_port,
		stop_server: make(chan bool),
		waitGroup:   &sync.WaitGroup{}}
}

func (s *TcpServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	s.waitGroup.Add(1)
	defer s.waitGroup.Done()

	for {
		select {
		case <-s.stop_server:
			fmt.Println("disconnecting", conn.RemoteAddr())
			return
		default:
		}

		conn.SetDeadline(time.Now().Add(1e9))

		buf := make([]byte, 4096)
		if _, err := conn.Read(buf); nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			fmt.Println(err)
			return
		}

		if _, err := conn.Write(buf); nil != err {
			fmt.Println(err)
			return
		}
	}
}

func (s *TcpServer) Listen() (err error) {
	s.waitGroup.Add(1)
	defer s.waitGroup.Done()

	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+strconv.Itoa(int(s.ListenPort)))
	if err != nil {
		fmt.Println("Started Listener fail.", err.Error())
	}

	tcplistener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Started Listener fail.", err.Error())
	}
	fmt.Println("Listener Started...", "Port:", tcpAddr.Port)
	defer tcplistener.Close()

	isStop := false
	for {

		select {
		case <-s.stop_server:
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

	return
}

func (s *TcpServer) Start() (err error) {
	fmt.Println("TcpServer", s.Name, "started.")
	go s.Listen()
	return nil
}

func (s *TcpServer) Stop() (err error) {
	close(s.stop_server)
	s.waitGroup.Wait()

	fmt.Println("TcpServer", s.Name, "stoped.")
	return nil
}
