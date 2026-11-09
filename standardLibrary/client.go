package main

import (
	"net"
)

func main() {
	// 建立连接
	conn, err := net.Dial("tcp", "0.0.0.0:1234")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 发送数据
	for i := range 10 {
		_, err := conn.Write([]byte{'a' + byte(i)})
		if err != nil {
			panic(err)
		}
	}
}
