package main

import (
	"gf-app/app/service"
	_ "gf-app/boot"
	_ "gf-app/router"
	"github.com/gogf/gf/frame/g"
)

func main() {

	go service.Tcpclient()

	g.Server().Run()
}
