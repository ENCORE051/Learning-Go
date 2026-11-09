package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

func main() {
	// 监听地址
	listener, err := net.Listen("tcp", "0.0.0.0:1234")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	var wg sync.WaitGroup

	for {
		// 阻塞等待下一个连接建立
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		// 开启一个新的协程去异步处理该连接
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := make([]byte, 4096)
			for {
				// 从连接中读取数据
				n, err := conn.Read(buf)
				if errors.Is(err, io.EOF) {
					break
				} else if err != nil {
					panic(err)
				}

				data := string(buf[:n])
				fmt.Println(data)
			}
		}()
	}

	wg.Wait()
}
