package main

import (
	"gf-app/app/service"
	"gf-app/app/util"
	_ "gf-app/boot"
	_ "gf-app/router"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

//动态ip端 用tcp server
func main() {

	service.Registe()

	go service.TcpServer()

	s := g.Server()
	s.SetPort(util.GetPort_client())

	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("nginx自动刷新公网ip服务正常！ 动态公网ip端")
	})
	s.Run()
}
