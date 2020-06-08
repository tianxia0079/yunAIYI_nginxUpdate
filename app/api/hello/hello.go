package hello

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

// Hello is a demonstration route handler for output "Hello World!".
func Hello(r *ghttp.Request) {
	r.Response.Writeln("nginx自动刷新电信公网ip服务正常!")
}

//云服务器使用
func UpdateIP(r *ghttp.Request) {
	ip := r.GetString("nowip")
	g.Log().Line(true).Debug("主动推送来 ip:", ip)
	//todo 1.更新tcp ip  2.更新nginx

	r.Response.WriteJson(g.Map{
		"state": "success",
	})
}
