package main

import (
	// "github.com/JackLeeMing/CloudNative/mpc"
	"os"

	// "github.com/JackLeeMing/CloudNative/server"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(6)
}

func main() {
	// 多生产者和多消费者
	// mpc.MPCExecute()
	// HTTP 服务
	// server.ExecuteServer()
	contextLogger := log.WithFields(log.Fields{
		"common": "this is a common field",
		"other":  "I also should be logged always",
	})
	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Info("Me too")
}
