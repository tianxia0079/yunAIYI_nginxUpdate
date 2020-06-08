package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gtcp"
)

//tcp长连接服务器端
func main() {
	fmt.Println("tcp服务启动...")
	gtcp.NewServer(":8999", func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				g.Log().Line(true).Debug("获取到信息:", string(data))
				//监听到有客户端输入信息
				if err := conn.Send(append([]byte("server "), data...)); err != nil {
					g.Log().Line(true).Debug("有客户端发送来消息但是回话失败", err)
				}
			}
			if err != nil {
				g.Log().Line(true).Error("客户端断开链接", err)
				break
			}
		}
	}).Run()
}
