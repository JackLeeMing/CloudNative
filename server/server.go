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
	"k8s.io/klog/v2"
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

// 3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	glog.V(4).Info("---- 进入 rootHandler ---")
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	delay := randInt(10, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	io.WriteString(w, "=================== Details of the http request header: ============\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
	glog.V(4).Infof("Respond in %d ms", delay)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	delay := randInt(10, 20)
	time.Sleep(time.Millisecond * time.Duration(delay))
	serviceFlag := os.Getenv("service_flag")
	urlPath := ""
	if serviceFlag == "service0" {
		urlPath = "service1"
	} else if serviceFlag == "service1" {
		urlPath = "service2"
	}
	if urlPath == "service1" || urlPath == "service2" {
		req, err := http.NewRequest("GET", "http://"+urlPath, nil)
		if err != nil {
			fmt.Printf("%s", err)
		}
		lowerCaseHeader := make(http.Header)
		for key, value := range r.Header {
			lowerCaseHeader[strings.ToLower(key)] = value
		}
		glog.Info("headers:", lowerCaseHeader)
		req.Header = lowerCaseHeader
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			glog.Info("HTTP get failed with error: ", "error", err)
			klog.Errorf("Failed to shutdown test server clearly: %v", err)
		} else {
			glog.Info("HTTP get succeeded")
		}
		if resp != nil {
			resp.Write(w)
		}
	} else {
		io.WriteString(w, "===================Details of the http request header:============\n")
		for k, v := range r.Header {
			io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
		}
	}
	glog.V(4).Infof("Respond in %d ms", delay)
}

func ExecuteServer() {
	level := os.Getenv("level")
	if level == "" {
		flag.Set("v", "4")
	} else {
		flag.Set("v", level)
	}
	metrics.Register()

	glog.V(2).Info("Starting http server...")
	httpport := os.Getenv("httpport")
	if httpport == "" {
		httpport = "80"
	}
	glog.V(4).Info("Server started and listing Port " + httpport + ".")
	mux := http.NewServeMux()
	mux.HandleFunc("/send", rootHandler)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/request1", request1Handler)
	mux.HandleFunc("/request2", request2Handler)
	mux.HandleFunc("/request3", request3Handler)
	mux.HandleFunc("/healthz", healthzHandler)

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/", healthzHandler)
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

func TracingServer() {
	level := os.Getenv("level")
	if level == "" {
		flag.Set("v", "4")
	} else {
		flag.Set("v", level)
	}
	metrics.Register()

	glog.V(2).Info("Starting http server...")
	httpport := os.Getenv("httpport")
	if httpport == "" {
		httpport = "80"
	}
	serviceFlag := os.Getenv("service_flag")
	glog.Infof("the service_flag is: " + serviceFlag)
	glog.V(4).Info("Server started and listing Port " + httpport + ".")
	mux := http.NewServeMux()
	mux.HandleFunc("/send", rootHandler)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/", homeHandler)
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
