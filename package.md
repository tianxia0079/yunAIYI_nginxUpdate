gf pack public,template boot/data.go --name boot
gf build main.go    

打包后及时删除data.go  否则开发环境资源文件优先级没有data.go高.


gf build mainServer.go -n Aliyun_nginxServer

gf build mainServer.go -n Aliyun_nginxServer
