package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"time"
)

const (
	retryMAX = 5
)

/**
长时间无法retry连接告警
*/
func main() {
	getnewip := "111.113.215.80"

	connectOK := false
	retry := 0
start:
	// Client
	c, err := gtcp.NewConn(getnewip + ":8999")
	if err != nil {
		retry++
		g.Log().Line(true).Error("和服务器链接失败了", err)

		//重试告警
		if retry == retryMAX {
			fmt.Println("重试" + gconv.String(retryMAX) + "次后失败请及时处理！")

		}
	} else {
		connectOK = true
		retry = 0
		//链接成功通知 第一次也通知，重试后也通知
		fmt.Println("链接成功！")
	}

	for {
		time.Sleep(1 * time.Second)

		if connectOK {
			if err := c.Send([]byte("1")); err != nil {
				//成功发出去了
				connectOK = false
				glog.Error("和服务器链接失败了", err)
			} else {
				data, err := c.Recv(-1)
				if len(data) > 0 {
					g.Log().Line(true).Debug("获取到服务器响应信息:", string(data))
				}
				if err != nil {
					//没有得到响应
					connectOK = false
					g.Log().Line(true).Error("客户端断开链接", err)
					break
				}
			}
		} else {
			goto start
		}

	}
}
