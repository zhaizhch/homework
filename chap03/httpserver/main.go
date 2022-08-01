package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
)

/*
课后练习 2.2
• 编写一个 HTTP 服务器，此练习为正式作业需要提交并批改
• 鼓励群里讨论，但不建议学习委员和课代表发满分答案，给大家留一点思考空间
• 大家视个人不同情况决定完成到哪个环节，但尽量把1都做完
1.接收客户端 request，并将 request 中带的 header 写入 response header
2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4.当访问 localhost/healthz 时，应返回200
*/
const (
	Version = "VERSION"
)

func main() {
	print("hello world")
	log.Println("hello hello")
	mux := http.NewServeMux()
	mux.HandleFunc("/readHeader", readHeader)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/recordLog", recordLog)
	mux.HandleFunc("/getEnv", getEnv)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	err := http.ListenAndServe(":80", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func readHeader(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
		log.Println("header info: ", fmt.Sprintf("%s=%s\n", k, v))
	}
}

func getEnv(w http.ResponseWriter, r *http.Request) {
	//VERSION写入环境变量
	os.Setenv(Version, "1.0.0")
	//读取env中的VERSION
	version := os.Getenv(Version)
	//写入response header
	io.WriteString(w, fmt.Sprintf("%s=%s/n", Version, version))
	w.Header().Set(Version, version)
	log.Println(Version, version)
}

func recordLog(w http.ResponseWriter, r *http.Request) {
	ip := ClientIP(r)
	log.Println("code", 200)
	log.Println("ip", ip)
	io.WriteString(w, fmt.Sprintf("ip:%s, code:%s/n", ip, "200"))
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "200")
	log.Println("healthz check!")
}

func getCurrentIP(r *http.Request) string {
	// 这里也可以通过X-Forwarded-For请求头的第一个值作为用户的ip
	//但是要注意的是这两个请求头代表的ip都有可能是伪造的
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		// 当请求头不存在即不存在代理时直接获取ip
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	return ip
}

// ClientIP 尽最大努力实现获取客户端 IP 的算法。
//解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
