// TcpServer project main.go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"tcpserver"
)

func wait_ctrl_c(tcp_server *tcpserver.TcpServer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for sig := range c {
		fmt.Printf("received ctrl+c (%v)\n", sig)

		tcp_server.Stop()
		fmt.Println("tcp server app stoped.")
		os.Exit(0)
	}
}

func main() {
	fmt.Println("tcp server app started.")

	tcp_server := tcpserver.NewTcpServer("virtual_tcp_server", 8888)
	tcp_server.Start()

	wait_ctrl_c(tcp_server)
}
