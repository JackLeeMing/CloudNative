package pc

import (
	"fmt"
	"math/rand"
	"time"
)

// 箭头指向 通道变量 发数据到通道
// 箭头指向 变量 从通道接收数据
func producer(tag string, ch chan<- string) {
	for {
		ch <- fmt.Sprintf("%s: %v", tag, rand.Int31())
		time.Sleep(time.Second)
	}
}

func consumer(ch <-chan string) {
	for {
		msg := <-ch
		fmt.Println(msg)
	}
}

func PCTest() {
	ch := make(chan string)
	go producer("cat", ch)
	go producer("dog", ch)

	consumer(ch)
}
