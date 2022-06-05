package main

import (
	"github.com/JackLeeMing/CloudNative/mpc"
	"github.com/JackLeeMing/CloudNative/server"
)

func main() {
	// 多生产者和多消费者
	mpc.MPCExecute()
	// HTTP 服务
	server.ExecuteServer()
}
