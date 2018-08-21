package main

import (
	"net/http"
	"fmt"
	"sync"
	"crypto/tls"
	"io/ioutil"
)

func main(){
	/*
		http://wxdd198a901fa24220.h5.inside.xiaoe-tech.com/use_invite_code
	*/
	url := "http://%s.h5.xiaoe-tech.com/use_invite_code?bizData[invite_code]=%s"

	//appId := "apppcHqlTPT3482"
	wxAppId := "wxdd198a901fa24220"
	batchId := "9180674321438155"

	url = fmt.Sprintf(url, wxAppId, batchId)
	fmt.Printf("url=%s", url)

	var workResultLock sync.WaitGroup

	loop := 1
	for {
		if loop >= 10 {
			break
		}
		loop++
		go func(requestUrl interface{}) {
			workResultLock.Add(1)

			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			client := &http.Client{}
			reqest, err := http.NewRequest("GET", url, nil)
			reqest.Header.Add("Cookie", "tgw_l7_route=d5c650a97a16252e1c20712f697ef070; xiaoe_session=eyJpdiI6ImNDdlNzc0llbnNoN2tEVEdSY2JscUE9PSIsInZhbHVlIjoiKzRZWDhqSmxNZFdBU3lHU3pFZ2xcL2RVTm40bWhWNnVOR0YrOE5ZbjhMSU1ZbWpnUnl0OGZkY3A1QVdhZVBzUVlveTVFdzNhampKTjY1aXhhNGw3bm13PT0iLCJtYWMiOiIzOTc4N2Q2NTU2MWJkYzk4ZGY4Mjc3YzJlZTkwN2MxZjNmYWFmMmRjNDI2YzE2MjMzYzNjODVjNzgzOWY4ZmU3In0%3D; Hm_lvt_17bc0e24e08f56c0c13a512a76c458fb=1532329464; Hm_lpvt_17bc0e24e08f56c0c13a512a76c458fb=1532852007")
			reqest.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1 wechatdevtools/1.02.1806120 MicroMessenger/6.5.7 Language/zh_CN webview/15323291605963160 webdebugger port/55334");

			response, err := client.Do(reqest)
			defer response.Body.Close()

			if err!=nil{
				fmt.Println("request has error,err=%v", err)
			}

			body, _:=ioutil.ReadAll(response.Body)
			fmt.Println("body:%v", string(body))
			fmt.Println("resp:%v", response)
			workResultLock.Done()
		}(url)
	}

	fmt.Println("done")
	workResultLock.Wait()
}
