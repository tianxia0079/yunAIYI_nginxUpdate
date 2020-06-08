package service

import (
	"fmt"
	"gf-app/app/util"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gmutex"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"time"
)

//云服务器
var (
	retryMAX  = 0
	DynamicIP = "111.1.1.1" //会被接口更新
)

func init() {
	fmt.Println("参数初始化...")
	retryMAX = util.Get_max_retry()
	if retryMAX <= 5 {
		panic("重试次数告警必须大于5")
	}
}

func UpdateIp(ip string) {
	lock := gmutex.New()
	lock.Lock()

	DynamicIP = ip

	confs := util.Get_nginx_confgs()

	need_reload_nginx := false;
	for _, v := range confs {
		v = gstr.Trim(v)
		newinfo := ""

		same := false;
		gfile.ReadLines(v, func(text string) {
			//fmt.Println("->" + text)
			if (gstr.Contains(text, "set $myip")) {
				if (gstr.Contains(text, ip)) {
					//相同ip 跳过
					same = true;
					return;
				} else {
					text = "set $myip \"" + ip + "\";"
				}
			}
			newinfo = newinfo + text + util.Get_huanhang();
		})
		if !same {
			gfile.PutContents(v, newinfo)
			g.Log().Line().Info("更新内容成功", newinfo)

			need_reload_nginx = true;
		}
	}

	if need_reload_nginx {
		shellerror := gproc.ShellRun("systemctl reload nginx")
		if shellerror != nil {
			g.Log().Line().Error("刷新nginx错误:", shellerror)
		} else {
			g.Log().Line().Info("nginx重启成功")
		}
	}

	lock.Unlock()

}
func Tcpclient() {

	connectOK := false
	retry := 0
start:
	// Client
	fmt.Println("DynamicIP " + DynamicIP)
	c, err := gtcp.NewConn(DynamicIP + ":8999")
	if err != nil {
		retry++
		g.Log().Line(true).Error("和服务器链接失败了", err)

		//重试告警
		if retry%retryMAX == 0 {
			fmt.Println("重试" + gconv.String(retry) + "次后失败请及时处理！")

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
				//没有成功发出去
				connectOK = false
				glog.Error("和服务器链接失败了", err)
			} else {
				data, err := c.Recv(-1)
				if len(data) > 0 {
					g.Log().Line(true).Debug("获取到服务器响应信息:", string(data))
				}
				if err != nil {
					//发出去了 没有得到响应
					connectOK = false
					g.Log().Line(true).Error("客户端断开链接", err)
					goto start
				}
			}
		} else {
			goto start
		}

	}
}
