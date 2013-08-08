package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

var ProxyApi = "http://thinklong.sinaapp.com/proxy.php?act=get&limit="
var ProxyDelApi = "http://thinklong.sinaapp.com/proxy.php?act=del&url="
var ProxyLogApi = "http://thinklong.sinaapp.com/executelog.php?host="

//安卓市场
//var Url = "http://apk.hiapk.com/Download.aspx?aid=1543012&rel=nofollow&module=256&info=41OLiFBORVw%3D"
//var Referer_url = "http://apk.hiapk.com/html/2013/06/1543012.html?module=256&info=41OLiFBORVw%3D"
//var Host = "apk.hiapk.com"

//应用宝
var Url = "http://android.myapp.com/android/down.jsp?appid=796684&icfa=-1&lmid=1022&g_f=0&actiondetail=0&softname=%E5%8F%A3%E8%A2%8B%E4%B9%90%E5%B1%85&downtype=1&transactionid=1375434227223079&topicid=-1&pkgid=-1"
var Referer_url = "http://android.myapp.com/android/appdetail.jsp?appid=796684&actiondetail=0&pageNo=1&clickpos=1&softname=%E5%8F%A3%E8%A2%8B%E4%B9%90%E5%B1%85&transactionid=1375434227223079&lmid=1022"
var Host = "android.myapp.com"

//var Url = "http://www.baidu.com/s?wd=ip"
//var Host = "www.baidu.com"
//var Referer_url = "http://www.baidu.com/s?wd=ip"

//var proxy_array []string
type http_pool struct {
	Url string
	//Referer_url  string
	Headers  map[string]string
	Method   string
	ProxyUrl string
	Timeout  int
	//	Sleep        int
}
type ProxyList struct {
	Id   string
	Url  string
	Host string
	Port string
	Time string
}

var ProxyData []json.RawMessage
var proxy_list ProxyList

func main() {
	count := 2000
	limit := fmt.Sprintf("%d", count)
	ProxyApi += limit
	//获取Proxy数据
	httpx := Newhttppool(ProxyApi, "", "GET", "", 10)

	proxy_data := httpx.http_send(1)
	//proxy_array1 := make(map[string]string{})
	//var proxy_array1 ProxyList

	if err := json.Unmarshal([]byte(proxy_data), &ProxyData); err != nil {
		log.Fatal(err.Error())
	}

	//proxy_array := make(map[int]string{})
	for key, _ := range ProxyData {
		json.Unmarshal(ProxyData[key], &proxy_list)

		httpx := Newhttppool(Url, Referer_url, "GET", string(proxy_list.Url), 10)

		httpx.http_send(0)
		//fmt.Println(i)
		time.Sleep(time.Duration(10) * time.Second)
	}
	//proxy_array = []string{"http://125.39.66.131:80/", "http://119.161.133.188:9000/"}
	/*
		for i := 0; i <= count; i++ {
			httpx := Newhttppool(Url, Referer_url, "GET", proxy_array[i], 3)
			httpx.http_send(0)
			time.Sleep(time.Duration(10) * time.Second)
		}
	*/
}
func Newhttppool(url string, referer string, method string, proxyurl string, timeout int) (h *http_pool) {
	headers := make(map[string]string)
	headers["Referer"] = referer
	headers["Host"] = Host
	//headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
	//headers["Accept-Encoding"] = "gzip,deflate,sdch"
	//headers["Accept-Language"] = "zh-CN,zh;q=0.8"
	//headers["Connection"] = "keep-alive"
	headers["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.116 Safari/537.36"
	return &http_pool{
		Url:      url,
		Headers:  headers,
		Method:   method,
		ProxyUrl: proxyurl,
		Timeout:  timeout,
	}
}
func del_proxy(id string) {
	url_ := ProxyDelApi + id
	fmt.Println(url_)
	httpx := Newhttppool(url_, "", "GET", "", 10)
	httpx.http_send(1)
}
func execute_log(host string) {
	log_url := ProxyLogApi + host
	httpx := Newhttppool(log_url, "", "GET", "", 5)
	httpx.http_send(1)
}

func (h *http_pool) http_send(is_content int) (html string) {

	transport := &http.Transport{}
	if len(h.ProxyUrl) > 0 {
		transport = getTransportFieldURL(h.ProxyUrl)
	}
	//设置超时时间
	//timeout
	transport.Dial = func(netw, addr string) (net.Conn, error) {
		deadline := time.Now().Add(2000 * time.Millisecond)
		c, err := net.DialTimeout(netw, addr, time.Second)
		if err != nil {
			return nil, err
		}
		c.SetDeadline(deadline)
		return c, nil
	}

	client := &http.Client{Transport: transport}
	req, err := http.NewRequest(h.Method, h.Url, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	//headers
	if len(h.Headers) > 0 {
		for k, v := range h.Headers {
			req.Header.Set(k, v)
		}
	}
	resq, err := client.Do(req)
	//defer resq.Body.Close()
	if err != nil {
		//log.Fatal(err.Error())
		//log.Println(err.Error())
		fmt.Println("请求失败！")
		html = "false"
		if len(h.ProxyUrl) > 0 {
			del_proxy(string(proxy_list.Id))
			fmt.Println(string(proxy_list.Id))
		}
		//resq.Body.Close()
		return
	}
	if resq.StatusCode == 200 {
		//html, _ = string(resq.Location())
		if is_content == 1 {
			robots, err := ioutil.ReadAll(resq.Body)
			if err != nil {
				log.Fatal(err.Error())
			}
			html = string(robots)
		} else {
			execute_log(h.Headers["Host"])
			html += "true\n"
			fmt.Println(html)
		}
		resq.Body.Close()
		//
	} else {
		html = ""
	}
	return
}

//指定代理ip
func getTransportFieldURL(proxyurl string) (transport *http.Transport) {
	url_i := url.URL{}
	url_proxy, _ := url_i.Parse(proxyurl)
	transport = &http.Transport{Proxy: http.ProxyURL(url_proxy)}
	return
}
