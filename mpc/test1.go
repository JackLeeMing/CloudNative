package mpc

import (
	"fmt"
	"sync"
	"time"
)

func GroupTest() {
	group := sync.WaitGroup{}
	group.Add(1)
	go func() {
		time.Sleep(time.Second)
		fmt.Println("Done.")
		group.Done()
	}()

	group.Wait()

	group.Add(1)
	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println("Done2.")
		group.Done()
	}()

	group.Wait()

	test2()
}
