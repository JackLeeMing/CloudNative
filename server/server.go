package server

import (
	"context"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/JackLeeMing/CloudNative/metrics"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/**
读锁 不互斥
写锁 互斥
*/

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

// 1. 接收客户端 request，并将 request 中带的 header 写入 response header
func request1Handler(response http.ResponseWriter, request *http.Request) {
	glog.V(4).Info("request1 handler")
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
	glog.V(4).Info("request from: ", request.RemoteAddr)
	verStr := os.Getenv("VERSION")
	response.Header().Set("VERSION", verStr)
	io.WriteString(response, fmt.Sprintf("%s=%s\n", "VERSION", verStr))
}

//3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
func request3Handler(response http.ResponseWriter, request *http.Request) {
	from := request.RemoteAddr
	ipStr := strings.Split(from, ":")
	glog.V(4).Info("Client-> ip= " + ipStr[0])
	glog.V(4).Info("Server-> response code= " + strconv.Itoa(int(http.StatusOK)))
	io.WriteString(response, "ok\n")
}

// 4.当访问 localhost/healthz 时，应返回 200
func healthzHandler(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(200) // http.StatusOK
	io.WriteString(response, "ok\n")
}

func rootHandler(response http.ResponseWriter, _ *http.Request) {
	glog.V(4).Info("---- 进入 rootHandler ---")
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	delay := randInt(10, 2000)
	verStr := os.Getenv("VERSION")
	logLevel := os.Getenv("loglevel")
	httpport := os.Getenv("httpport")
	values := []string{verStr, logLevel, httpport}

	time.Sleep(time.Millisecond * time.Duration(delay))
	io.WriteString(response, strings.Join(values, ","))
	glog.V(4).Info("rootHandler 在%dms内 完成响应", delay)
}

func ExecuteServer() {
	level := os.Getenv("level")
	if level == "" {
		flag.Set("v", "4")
	} else {
		flag.Set("v", level)
	}
	glog.V(2).Info("Starting http server...")
	httpport := os.Getenv("httpport")
	if httpport == "" {
		httpport = "8090"
	}
	metrics.Register()
	glog.V(4).Info("Server started and listing Port " + httpport + ".")
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/request1", request1Handler)
	mux.HandleFunc("/request2", request2Handler)
	mux.HandleFunc("/request3", request3Handler)
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/send", rootHandler)

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	srv := http.Server{
		Addr:    ":" + httpport,
		Handler: mux,
	}
	// golang 优雅终止
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glog.Fatalf("listen: %s\n", err)
		}
	}()
	glog.Infof("--- Server started ---")
	<-done
	glog.Infof("--- Server stopped ---")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		glog.Fatalf("Server Shutdown Failed: %+v", err)
	}
	glog.Info("Server Exited Properly")
}
