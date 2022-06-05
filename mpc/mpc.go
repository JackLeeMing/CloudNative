package mpc

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var timeOut = false

func producer(threadId int, wg *sync.WaitGroup, ch chan string) {
	count := 0
	for !timeOut {
		time.Sleep(1 * time.Second)
		count++
		data := strconv.Itoa(threadId) + "--" + strconv.Itoa(count)
		fmt.Printf("Producer. %s\n", data)
		ch <- data
	}
	wg.Done()
}

func consumer(wg *sync.WaitGroup, ch chan string) {
	for data := range ch {
		time.Sleep(1 * time.Second)
		fmt.Printf("Consumer. %s\n", data)
	}
	wg.Done()
}

func MPCExecute() {
	ch := make(chan string, 10)
	wgP := new(sync.WaitGroup)
	wgC := new(sync.WaitGroup)
	for i := 0; i < 3; i++ {
		wgP.Add(1)
		go producer(i, wgP, ch)
	}
	for i := 0; i < 3; i++ {
		wgC.Add(1)
		go consumer(wgC, ch)
	}
	go func() {
		time.Sleep(3 * time.Second)
		timeOut = true
		fmt.Println("超时")
	}()
	wgP.Wait()
	close(ch)
	wgC.Wait()
}
