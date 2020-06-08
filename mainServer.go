package main

import (
	"gf-app/app/service"
	"gf-app/app/util"
	_ "gf-app/boot"
	_ "gf-app/router"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
)

//云服务器端 用tcp client
func main() {
	//判断文件是否存在
	confs := util.Get_nginx_confgs()
	for _, v := range confs {
		v = gstr.Trim(v)

		if (!gfile.Exists(v)) {
			panic("nginx配置文件不存在,路径:" + v)
		}
	}

	go service.Tcpclient()

	s := g.Server()
	s.SetPort(util.GetPort_server())
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("nginx自动刷新公网ip服务正常！ 固定ip云服务器端")
	})
	s.BindHandler("/UpdateIP", func(r *ghttp.Request) {
		ip := r.GetString("nowip")
		timestamp := r.GetString("timestamp")
		sign := r.GetString("sign")
		g.Log().Line(true).Debug("主动推送来 ip:", ip)

		// 验证签名 	sign, _ := gmd5.EncryptString(ip + timestamp + sign_key)

		mysigntemp, _ := gmd5.EncryptString(ip + timestamp + service.Sign_key)
		state := ""

		if mysigntemp == sign {
			state = "success"

			// 1.更新tcp ip  2.更新nginx

			service.UpdateIp(ip)
		} else {
			state = "fail"
		}

		r.Response.WriteJson(g.Map{
			"state": state,
		})
	})
	s.Run()
}
