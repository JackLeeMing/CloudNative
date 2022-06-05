package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"
	"strings"
)

/**
读锁 不互斥
写锁 互斥

*/

// 1. 接收客户端 request，并将 request 中带的 header 写入 response header
func request1Handler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("request1 handler")
	headers := request.Header
	// type Header map[string][]string
	for header := range headers {
		values := headers[header]
		for i, v := range values {
			values[i] = strings.TrimSpace(v)
		}
		// 写入 header中
		response.Header().Set(header, strings.Join(values, ","))
	}
	io.WriteString(response, "ok\n")
}

// 2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
func request2Handler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("request from: ", request.RemoteAddr)
	verStr := os.Getenv("VERSION")
	response.Header().Set("VERSION", verStr)
	io.WriteString(response, fmt.Sprintf("%s=%s\n", "VERSION", verStr))
}

//3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
func request3Handler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("request from: ", request.RemoteAddr)
	from := request.RemoteAddr
	println("Client-> ip:port= " + from)
	ipStr := strings.Split(from, ":")
	println("Client-> ip= " + ipStr[0])
	println("Server-> response code= " + strconv.Itoa(int(http.StatusOK)))
	io.WriteString(response, "ok\n")
}

// 4.当访问 localhost/healthz 时，应返回 200
func healthzHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("request from: ", request.RemoteAddr)
	response.WriteHeader(200) // http.StatusOK
	io.WriteString(response, "ok\n")
}

func ExecuteServer() {
	os.Setenv("VERSION", "v1.0.1")
	fmt.Println("Server started and listing Port 8090.")
	mux := http.NewServeMux()
	mux.HandleFunc("/request1", request1Handler)
	mux.HandleFunc("/request2", request2Handler)
	mux.HandleFunc("/request3", request3Handler)
	mux.HandleFunc("/healthz", healthzHandler)

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	err := http.ListenAndServe(":8090", mux)
	if err != nil {
		log.Fatal(err)
	}
}
