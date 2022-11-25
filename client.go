package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("Dial err:", err)
		return
	}
	defer conn.Close() //关闭双向连接
	for {
		var inputMsg string
		fmt.Println("请输入你要发送的信息:")
		fmt.Scanln(&inputMsg)
		inputMsg = strings.Trim(inputMsg, "\r\n")
		if strings.ToUpper(inputMsg) == "quit" {
			return
		}
		_, err = conn.Write([]byte(inputMsg))
		if err != nil {
			fmt.Println("send err:", err)
			return
		}
		buf := [512]byte{}
		serverMsg, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("recv failed err:", err)
			return
		}
		fmt.Println("server message:", string(buf[:serverMsg]))
	}

}
