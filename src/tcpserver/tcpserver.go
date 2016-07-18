package tcpserver

import (
	//	"Net"
	"fmt"
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

func (s *TcpServer) Listen(stop_listener, listener_is_stopped chan bool) (err error) {
	/*
		listener, err := net.Listen("tcp", ":8888")
		if err != nil {
			fmt.Println("listen error:", err)
			return
		}
		listener.SetDeadline(time.now().Add(1*time.Second))
	*/
	fmt.Println("Listener Started.")

	/*
	   	for {
	   		c, err := listener.Accept()
	   		if err != nil {
	   			fmt.Println("accept error:", err)
	   			break
	   		}
	   		// start a new goroutine to handle
	   		// the new connection.
	      		go handleConn(c)
	   	}
	*/
	_ = <-s.stop_listener
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
