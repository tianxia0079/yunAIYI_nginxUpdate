package util

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"strconv"
)

// 打开系统默认浏览器

var commands = map[string]string{
	"windows": "start",
	"darwin":  "open",
	"linux":   "xdg-open",
}

var Phone_agent = []string{"iphone", "android", "ipad", "phone", "mobile", "wap", "netfront", "java", "opera mobi",
	"opera mini", "ucweb", "windows ce", "symbian", "series", "webos", "sony", "blackberry", "dopod",
	"nokia", "samsung", "palmsource", "xda", "pieplus", "meizu", "midp", "cldc", "motorola", "foma",
	"docomo", "up.browser", "up.link", "blazer", "helio", "hosin", "huawei", "novarra", "coolpad", "webos",
	"techfaith", "palmsource", "alcatel", "amoi", "ktouch", "nexian", "ericsson", "philips", "sagem",
	"wellcom", "bunjalloo", "maui", "smartphone", "iemobile", "spice", "bird", "zte-", "longcos",
	"pantech", "gionee", "portalmmm", "jig browser", "hiptop", "benq", "haier", "^lct", "320x320",
	"240x320", "176x220", "w3c ", "acs-", "alav", "alca", "amoi", "audi", "avan", "benq", "bird", "blac",
	"blaz", "brew", "cell", "cldc", "cmd-", "dang", "doco", "eric", "hipt", "inno", "ipaq", "java", "jigs",
	"kddi", "keji", "leno", "lg-c", "lg-d", "lg-g", "lge-", "maui", "maxo", "midp", "mits", "mmef", "mobi",
	"mot-", "moto", "mwbp", "nec-", "newt", "noki", "oper", "palm", "pana", "pant", "phil", "play", "port",
	"prox", "qwap", "sage", "sams", "sany", "sch-", "sec-", "send", "seri", "sgh-", "shar", "sie-", "siem",
	"smal", "smar", "sony", "sph-", "symb", "t-mo", "teli", "tim-", "tosh", "tsm-", "upg1", "upsi", "vk-v",
	"voda", "wap-", "wapa", "wapi", "wapp", "wapr", "webc", "winw", "winw", "xda", "xda-",
	"Googlebot-Mobile"}

func GetPort_client() int {
	httpport := g.Cfg().GetInt("server.AddressClient")
	return httpport
}

func GetPort_server() int {
	httpport := g.Cfg().GetInt("server.AddressServer")
	return httpport
}
func GetYunserverIP() string {
	httpport := g.Cfg().GetString("myconfig.my_yunserver_ip")
	return httpport
}

func Get_sign_key() string {
	sign_key := g.Cfg().GetString("myconfig.sign_key")
	return sign_key
}
func Get_huanhang() string {
	return "\r\n";

	system := runtime.GOOS
	switch system {
	case "windows":
		{
			return "\r\n"
		}
	case "linux":
		{
			return "\r"

		}
	default:
		return "\n"
	}
}
func Get_nginx_confgs() []string {
	return g.Cfg().GetStrings("myconfig.nginx_conf_paths")
}
func Get_max_retry() int {
	sign_key := g.Cfg().GetInt("myconfig.max_retry")
	return sign_key
}
func GetRedirectURL(url string) string {
	c := ghttp.NewClient()
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	response, _ := c.Get(url)
	return gconv.String(response.Header.Get("Location"))
}
func GetThisServerIP() string {
	ip := ghttp.GetContent("http://pv.sohu.com/cityjson?ie=utf-8")

	if gstr.Contains(ip, "403 Forbidden") {
		return GetThisServerIP()
	} else {
		ip = gstr.Trim(gstr.Split(gstr.Split(ip, "=")[1], ";")[0])

		jsonstr, _ := gjson.LoadContent(ip, true)
		return jsonstr.GetString("cip")
	}

}

//通过后端判断是否手机设备
func IsPhone(r *ghttp.Request) bool {
	agent := r.UserAgent()
	for _, v := range Phone_agent {
		if gstr.Pos(gstr.ToLower(agent), v) >= 0 {
			//fmt.Println("手机")
			return true
		}
	}
	return false
}

//通过后端判断是否 微信小程序
func Iswxminiapp(r *ghttp.Request) bool {
	agent := r.UserAgent()
	if gstr.Pos(gstr.ToLower(agent), gstr.ToLower("miniProgram")) >= 0 {
		//fmt.Println("手机")
		return true
	}
	return false
}

func HttpGet(apiUrl string) string {
	resp, err := http.Get(apiUrl)
	if err != nil {
		g.Log().Line(true).Error("http请求失败:", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return (string(body))
}

func GetMacAddrs() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return "fail to get mac"
	}
	macAddrs := []string{}
	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		macAddrs = append(macAddrs, macAddr)
	}
	return gconv.String(macAddrs)
}
func StaticFilePath_Root() string {
	abs := gfile.Pwd()
	abs = abs + gfile.Separator + "staicfile_root" + gfile.Separator
	if gfile.Exists(abs) {

	} else {
		gfile.Mkdir(abs)
	}
	return abs
}
func OpenUrl(uri string) {
	system := runtime.GOOS
	switch system {
	case "windows":
		{
			// 无GUI调用
			exec.Command(`cmd`, `/c`, `start`, uri).Start()
			break
		}
	case "linux":
		{
			exec.Command(`xdg-open`, uri).Start()
			break
		}
	case "darwin":
		{
			exec.Command(`open`, uri).Start()
			break
		}
	}
}

//2020-01-01 09:11:27
func GetbeijingTime() string {
	re := ""
	content := ghttp.GetContent("http://api.m.taobao.com/rest/api3.do?api=mtop.common.getTimestamp")
	taobaoTime := new(TaobaoTime)
	err := gjson.DecodeTo(content, &taobaoTime)
	if err != nil {
		g.Log().Line(true).Error(err)
		return ""
	}
	timestamptemp := gconv.Int64(taobaoTime.Data.T) / 1000
	re = gtime.NewFromTimeStamp(timestamptemp).String()
	g.Log().Line(true).Info(re)
	return re
}
func UpdateSystemDateAuto() {
	UpdateSystemDate(GetbeijingTime())

}

//2020-1-1 16:55:51 格式
func UpdateSystemDate(dateTime string) bool {
	system := runtime.GOOS
	switch system {
	case "windows":
		{
			_, err1 := gproc.ShellExec(`date  ` + gstr.Split(dateTime, " ")[0])
			_, err2 := gproc.ShellExec(`time  ` + gstr.Split(dateTime, " ")[1])
			if err1 != nil && err2 != nil {
				g.Log().Line(true).Info("更新系统时间错误:请用管理员身份启动程序!")
				return false
			}
			g.Log().Line(true).Info("已自动同步时间: ", dateTime)
			return true
			break
		}
	case "linux":
		{
			_, err1 := gproc.ShellExec(`date -s  "` + dateTime + `"`)
			if err1 != nil {
				g.Log().Line(true).Info("更新系统时间错误:", err1.Error())
				return false
			}
			g.Log().Line(true).Info("已自动同步时间: ", dateTime)
			return true
			break
		}
	case "darwin":
		{
			//todo:mac是否可以执行 未测试
			_, err1 := gproc.ShellExec(`date -s  "` + dateTime + `"`)
			if err1 != nil {
				g.Log().Line(true).Info("更新系统时间错误:", err1.Error())
				return false
			}
			g.Log().Line(true).Info("已自动同步时间: ", dateTime)
			return true
			break
		}
	}
	return false
}
func GOID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
func GetGoruntinesCount() string {
	var preback *pprof.Profile
	for _, p := range pprof.Profiles() {
		if p.Name() == "goroutine" {
			preback = p
			break
		}
	}
	return "正在运行中的goruntine数量: " + gconv.String(preback.Count())
}

type TimeStamp struct {
	T string `json:t`
}
type TaobaoTime struct {
	Api  string    `json:api`
	V    string    `json:v`
	ret  string    `json:ret`
	Data TimeStamp `json:data`
}
