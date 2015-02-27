package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const (
	apiURL = "http://wifiapi02.51y5.net/wifiapi/fa.cmd"
)

func main() {
	data := url.Values{
		"st":     {"m"},                                // 固定字段
		"appid":  {"001"},                              // 固定字段
		"v":      {"508"},                              // 固定字段
		"och":    {"guanwang"},                         // 固定字段
		"chanid": {"guanwang"},                         // 固定字段
		"pid":    {"qryapwd:commonswitch"},             // 请求类别
		"method": {"getSecurityCheckSwitch"},           // 请求方法
		"uhid":   {"a0000000000000000000000000000001"}, // 固定字段
		"mac":    {"c0:61:18:44:8a:98"},                // 本机 MAC
		"sign":   {""},                                 // 签名
		"ssid":   {"TP-LINK_8A98"},                     // WIFI SSID
		"bssid":  {"c0:61:18:44:8a:98"},                // WIFI MAC
		"dhid":   {"40289ec14942672d014954ad909a1147"}, // 设备字段
	}

	// 获取 sign
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var m Msg
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	data["sign"] = []string{getSign(data, m.RetSn)}

	// 正式请求
	resp, err = http.PostForm(apiURL, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}

// sign 计算方法为：
// 将要发送的数据依照 map 的 key 冒泡排序后
// 将值链接，然后再链接 retSn 得到字符串
// 最后取字符串的大写的MD5值为 sign
func getSign(data url.Values, retSn string) string {
	var buf = make([]string, len(data))
	var i int
	for k := range data {
		buf[i] = k
		i++
	}
	sort.Strings(buf)

	var str string
	for _, v := range buf {
		for _, v := range data[v] {
			str += v
		}
	}
	str += retSn

	return getUpperMD5(str)
}

func getUpperMD5(str string) string {
	return strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(str))))
}

type Msg struct {
	RetSn string `json:"retSn"`
}
