package main

import (
	"fmt"
	"time"

	"github.com/golang/glog"
)

/**
读锁 不互斥
写锁 互斥
go mod tidy 扫描依赖
*/
type Config struct {
	Host string `json:"host"`
	IP   string `json:"ip"`
}

func main() {
	config := Config{Host: "localhost", IP: "127.0.0.1"}
	glog.Infof("Hello")
	for i := 0; i < 2; i++ {
		go func(index int) {
			config.Host = fmt.Sprintf("%s-%d", config.Host, index)
			fmt.Println(config)
		}(i)
	}
	time.Sleep(time.Second)
	fmt.Println(config)
}
