package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

//处理请求
func process(conn net.Conn) {
	defer conn.Close() // 关闭链接通道
	for {
		reader := bufio.NewReader(conn)
		var buf [1024]byte
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			fmt.Print("read form client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("client message:", recvStr)
		var inputMsg string
		fmt.Println("请输入你要发送的信息:")
		fmt.Scanln(&inputMsg)
		inputMsg = strings.Trim(inputMsg, "\r\n")
		conn.Write([]byte(inputMsg))
	}
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		conn, err := listen.Accept() //建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn)
	}
}
