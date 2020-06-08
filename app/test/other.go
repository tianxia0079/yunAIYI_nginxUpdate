package main

import (
	"fmt"
	"gf-app/app/service"
	"io/ioutil"
	"net/http"
)

func main() {
	service.UpdateIp("1.1.1.1")
}

func login() string {
	url := "http://192.168.0.1/router/interface_status_show.cgi"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("获取地址错误")
	}
	req.Header.Set("Cookie", `__guid=224789850.346329641996463600.1583585485503.16; __DC_gid=224789850.565324479.1583585485504.1586937598450.36; Qihoo_360_login=831701798dddee1f49ff1a56d5ea6762`)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("登录错误")
	}
	resp_byte, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	respHtml := string(resp_byte)
	return respHtml
}
