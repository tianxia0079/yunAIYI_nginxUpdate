package service

import (
	"fmt"
	"gf-app/app/util"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"time"
)

var (
	httpapi_ip   = ""
	httpapi_port = ""
	Sign_key     = ""
)

func init() {
	httpapi_ip = util.GetYunserverIP()

	httpapi_port = gconv.String(util.GetPort_server())
	Sign_key = util.Get_sign_key()
}

//动态ip端

//客户端用
func GetNowIp() string {
	return util.GetThisServerIP()
}

func Registe() {
	g.Log().Line().Debug("开始注册最新公网ip....start")
	ip := util.GetThisServerIP()

	timestamp := gconv.String(gtime.Now().Timestamp())
	sign, _ := gmd5.EncryptString(ip + timestamp + Sign_key)
	url := "http://" + httpapi_ip + ":" + httpapi_port + "/UpdateIP?nowip=" + ip + "&timestamp=" + timestamp + "&sign=" + sign
	fmt.Println(url)
	content := ghttp.GetContent(url)
	g.Log().Line().Debug(content)
	if gstr.Contains(content, "success") {
		fmt.Println("已连接您的云服务器并更新成功!")
	} else {
		g.Log().Line().Error("注册失败，递归注册!")
		time.Sleep(1 * time.Second)
		Registe()
	}
}

func TcpServer() {
	//开启本地8999端口 tcp
	gtcp.NewServer(":8999", func(conn *gtcp.Conn) {
		defer conn.Close()

		remoteip := conn.RemoteAddr()
		fmt.Println(remoteip.String() + " remoteip")
		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				//g.Log().Line(true).Debug("获取到信息:", string(data))
				//监听到有客户端输入信息
				if err := conn.Send(append([]byte("server "), data...)); err != nil {
					g.Log().Line(true).Debug("有客户端发送来消息但是回话失败", err)
					Registe()
					break
				}
			}
			if err != nil {
				g.Log().Line(true).Error("客户端断开链接", err)

				Registe()
				break
			}
		}
	}).Run()
}
