package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"time"
)

const (
	corpid  = "ww7064a3ad71f6a148" //企业ID
	agentId = "1000002"            //应用ID
	secret  = ""                   //Secret
	url     = "https://www.binance.com/zh-CN/support/announcement/c-48"
)

func main() {
	title := getLatestCoinDeployNotice()
	for true {
		time.Sleep(60 * time.Second)
		newTitle := getLatestCoinDeployNotice()
		if newTitle != title {
			msg, err := SendCardMsg("wangtianyi", "binance", title, url)
			if err != nil {
				log.Println(msg)
			}
			title = newTitle
		}
	}
}

func getLatestCoinDeployNotice() string {
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	//out, err := os.Create("1.html")
	//defer out.Close()
	//_, err = io.Copy(out, res.Body)

	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
	}
	title := ""
	doc.Find("#__APP_DATA").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		//ioutil.WriteFile("1.json", []byte(text), 0644)

		jsonParsed, _ := gabs.ParseJSON([]byte(text))

		// Search JSON
		//fmt.Println(jsonParsed.Path("routeProps.b723.catalogs.0.articles").Data())
		for key, child := range jsonParsed.Search("routeProps", "b723", "catalogs", "0", "articles", "0").ChildrenMap() {
			if key == "title" {
				title = child.Data().(string)
				fmt.Printf("Key=>%v, Value=>%v\n", key, title)
				break
			}
		}
	})
	if res != nil && res.Body != nil {
		res.Body.Close()
	}
	return title
}

// 企业微信应用消息提醒方法如下
func SendCardMsg(ToUsers interface{}, title, description, url string) (map[string]interface{}, error) {
	btntxt := "详情" //可自定义卡片底下的文字

	qyurl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", corpid, secret)
	data, err := httpGetJson(qyurl)
	if err != nil {
		log.Println(err)
		return data, err
	}

	errcode := data["errcode"].(float64)
	if errcode != 0 {
		log.Println(errcode)
		return make(map[string]interface{}), nil
	}
	access_token := data["access_token"]
	//卡片内容，不同类型消息通知以下内容需修改(可参考企业微信开发文档)
	req := map[string]interface{}{
		"touser":  ToUsers,
		"msgtype": "textcard",
		"agentid": agentId,
		"textcard": map[string]interface{}{
			"title":       title,
			"description": description,
			"url":         url,
			"btntext":     btntxt,
		},
	}

	sendurl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", access_token)
	data, err = httpPostJson(sendurl, req)
	fmt.Println(data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, nil
}

// 封装了http请求的GET和POST方法，其返回值都是map[string]interface{}，方便我们使用。
func httpGetJson(url string) (map[string]interface{}, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func httpPostJson(url string, data map[string]interface{}) (map[string]interface{}, error) {
	res, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(res))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data2 map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data2, nil
}
