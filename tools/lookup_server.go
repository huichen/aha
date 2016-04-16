package main

import (
	"../json_service"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	host         = flag.String("host", "", "运行HTTP服务的IP")
	port         = flag.String("port", "9999", "运行HTTP服务的端口")
	logFile      = flag.String("log_file", "/tmp/aha_http.log", "HTTP log文件位置")
	staticFolder = flag.String("static", "../static", "HTTP静态文件位置")
	service      json_service.JsonService
)

func main() {
	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	// 捕获ctrl-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			log.Print("捕获Ctrl-c，退出服务器")
			os.Exit(0)
		}
	}()

	// 打开log文件
	file, err := os.Create(*logFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	go service.Init()

	http.HandleFunc("/metric", service.MetricJsonRpcService)
	http.HandleFunc("/tag", service.TagJsonRpcService)
	http.HandleFunc("/option", service.OptionJsonRpcService)
	http.HandleFunc("/stats", service.StatsService)
	http.Handle("/", http.FileServer(http.Dir(*staticFolder)))
	log.Print("服务器启动")
	http.ListenAndServe(*host+":"+*port, WriteLog(http.DefaultServeMux, file))
}

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	w.length = len(b)
	return w.ResponseWriter.Write(b)
}

func WriteLog(handle http.Handler, fileHandler *os.File) http.HandlerFunc {
	logger := log.New(fileHandler, "", 0)
	return func(w http.ResponseWriter, request *http.Request) {
		start := time.Now()
		writer := statusWriter{w, 0, 0}
		handle.ServeHTTP(&writer, request)
		end := time.Now()
		latency := end.Sub(start)
		statusCode := writer.status
		length := writer.length
		if request.URL.RawQuery != "" {
			logger.Printf("%v %s %s \"%s %s%s%s %s\" %d %d \"%s\" %v", end.Format("2006/01/02 15:04:05"), request.Host, request.RemoteAddr, request.Method, request.URL.Path, "?", request.URL.RawQuery, request.Proto, statusCode, length, request.Header.Get("User-Agent"), latency)
		} else {
			logger.Printf("%v %s %s \"%s %s %s\" %d %d \"%s\" %v", end.Format("2006/01/02 15:04:05"), request.Host, request.RemoteAddr, request.Method, request.URL.Path, request.Proto, statusCode, length, request.Header.Get("User-Agent"), latency)
		}
	}
}
