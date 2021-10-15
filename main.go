package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// 代理来源 https://www.89ip.cn/tqdl.html?api=1&address=美国&num=200

var (
	// 这里列出要获取代理的网站列表
	SITES = "https://www.89ip.cn/tqdl.html?api=1&address=美国&num=200"
	// 超时时间，单位：秒
	TIMEOUT = 5
	// 协程管理
	wg sync.WaitGroup
)

// 用代理IP访问指定的地址来测试一下是否
func ProxyTest(proxy_addr string) {
	defer wg.Done()
	// 检测代理IP访问地址
	var testUrl string
	// 判定传来的代理是否是https
	if strings.Contains(proxy_addr, "https") {
		testUrl = "https://httpbin.org/get"
	} else {
		testUrl = "http://httpbin.org/get"
	}
	// 解析代理地址
	proxy, err := url.Parse(proxy_addr)
	if err != nil {
		fmt.Println("代理地址发生错误：", err)
	}
	// 创建连接客户端
	httpClient := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			Proxy:                 http.ProxyURL(proxy),
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * time.Duration(TIMEOUT),
		},
	}
	// 用于请求时设置的参数
	request, err := http.NewRequest("GET", testUrl, nil)
	if err != nil {
		fmt.Println("参数错误")
	}
	// 添加头部协议
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-US;q=0.7,ko;q=0.6,zh-TW;q=0.5")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Mobile Safari/537.36")
	// 判断代理访问时间
	begin := time.Now()
	// 发起请求
	response, err := httpClient.Do(request)
	if err != nil {
		// log.Println(err)
		return
	}
	defer response.Body.Close()
	speed := int(time.Now().Sub(begin).Nanoseconds() / 1000 / 1000)
	body, err := ioutil.ReadAll(response.Body)
	// 这里验证一下是否真的可以代理成功
	if !strings.Contains(string(body), testUrl) {
		return
	}
	// 判断是否成功访问，如果成功访问StatusCode应该是200
	if response.StatusCode != http.StatusOK {
		// log.Println(err)
		return
	}
	if response.StatusCode == 200 {
		fmt.Println(proxy_addr, "ok", speed)
	} else {
		fmt.Println("err", err)
	}
}

func main() {
	fi, err := os.Open("proxy.txt")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(0)
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	for {
		cont, _, end := br.ReadLine()
		if end == io.EOF {
			break
		}
		wg.Add(1)
		go ProxyTest("http://" + string(cont))
	}
	wg.Wait()
}
