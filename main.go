package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// 定义代理服务器监听的地址
	proxyAddr := ":8080"

	// 启动代理服务器
	proxy := &Proxy{}
	log.Fatal(http.ListenAndServe(proxyAddr, proxy))
}

// 定义代理结构体
type Proxy struct{}

// 实现http.Handler接口的ServeHTTP方法
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 打印请求信息
	fmt.Printf("%s %s %s\n", r.Method, r.Host, r.URL.Path)

	// 创建一个新的HTTP请求对象
	req, err := http.NewRequest(r.Method, "https://api.openai.com"+r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 复制源请求的HTTP头部
	req.Header = r.Header

	// 发送HTTP请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 复制响应的HTTP头部
	for key, value := range resp.Header {
		w.Header().Set(key, value[0])
	}

	// 返回响应的HTTP状态码和响应体
	w.WriteHeader(resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
