package main

import (
	"fmt"
	"github.com/golang/glog"
	"io"
	"net"
	"net/http"
	"os"
)

/**
1. 使用 mux,err := http.ListenAndServe(":80", mux)
2. 直接使用 http.HandleFunc("/", rootHandler),直接
*/

func main() {
	glog.V(2).Info("Starting http server...")
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/writeToRepHeader", writeToRepHeader)
	mux.HandleFunc("/writeVersionToRepHeader", writeVersionToRepHeader)

	err := http.ListenAndServe(":80", logRequestHandler(mux))
	if err != nil {
		fmt.Println(err)
	}

}

/**
Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
*/
func logRequestHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// 在我们包装的时候调用原始的 http.Handle
		h.ServeHTTP(w, r)

		// 得到请求的有关信息，并记录之
		uri := r.URL.String()
		method := r.Method
		// ... 更多信息
		ip := GetIPFromRequest(r)
		glog.Infof("%s,%s,%s\n", uri, method, ip)
		//glog.InfoDepth(1, "This is info message", 12345)
	}

	// 用 http.HandlerFunc 包装函数，这样就实现了 http.Handler 接口
	return http.HandlerFunc(fn)
}

func GetIPFromRequest(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		glog.Errorf("userip: %q is not IP:port", r.RemoteAddr)
		return ""
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		glog.Errorf("userip: %q is not IP:port", r.RemoteAddr)
		return ""
	}
	return userIP.String()
}

//当访问 localhost/healthz 时，应返回200
func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "200\n")
}

//接收客户端 request，并将 request 中带的 header 写入 response header
func writeToRepHeader(w http.ResponseWriter, r *http.Request) {
	header := r.Header
	responseHeader := w.Header()
	for key, values := range header {
		for _, value := range values {
			responseHeader.Add(key, value)
		}
	}
	io.WriteString(w, "writeToRepHeader")
}

//读取当前系统的环境变量中的 VERSION 配置，并写入 response header
func writeVersionToRepHeader(w http.ResponseWriter, r *http.Request) {
	ver := os.Getenv("VERSION")
	responseHeader := w.Header()
	responseHeader.Add("VERSION", ver)
	io.WriteString(w, "writeToRepHeader")
}
