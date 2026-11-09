package main

import "fmt"

func main() {
	ch := make(chan int, 1)
	defer close(ch)
	go write(ch)
	fmt.Println(<-ch)
}

func write(ch chan<- int) {
	// 只能对管道发送数据
	ch <- 1
	close(ch)
}
