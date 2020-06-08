# yunAIYI_nginxUpdate
## 业务背景
使用第三方提供的内网穿透服务，还是有明显的网络延迟，如果有阿里云或腾讯云，可以自己结合宽带动态公网ip部署自己的内网穿透。没有过多的流量回路，速度肯定比第三方内网
       穿透快和稳定。

原理：TCP长连接状态处理，当宽带动态公网ip变化，实时更新到阿里云web服务器，更新nginx，从而完成内网穿透。

## 前置条件
1. 电信宽带下的服务器 PCserver1。有电信公网ip但动态变化，且无80 443端口权限。
2. 阿里云web服务器例 PCserver2，固定公网ip如：1.1.1.1 。部署web服务会nginx。 
(测试环境服务器配置：5m宽带，1核2g内存。)    
nginx配置，把80 443都代理到一个变量ip,内容随便写，会自动更新。  
set @mydynamicIP "111.1.1.1"。   
配置信息参考后面示例。 

3. 域名如 mydomain.top，a记录绑定到1.1.1.1。


## 最终效果
核心业务系统部署在宽带环境下的pc服务器上。（自己的pc服务器硬件性能远超这台云服务器）
访问 https://mydomain.top ，仅通过PCserver2 nginx进行代理，业务服务器还是在PCserver1上。


## 使用说明：
 Aliyun_nginxServer Aliyun_nginxClient 启动顺序可以随意 。     

1. 云服务器 PCserver2 下载运行  Aliyun_nginxServer     
双击运行即可。         
访问： 有内容，说明启动成功。

2. 自己的pc PCserver1 下载运行  Aliyun_nginxClient。
保障8999端口没有被占用。
配置文件中修改 aliyunip 为自己的阿里云web服务器ip。       
双击运行即可。     
访问： 有内容，说明启动成功。 


## nginx配置文件示例
```
#homeapp.top.conf

server {

		set $myip "111.113.215.87";

		listen 443 ssl http2;

		server_name homeapp.top;

		# SSL
		ssl_certificate /usr/local/nginx/mywebsiteconf/www.homeapp.top.pem;
		ssl_certificate_key /usr/local/nginx/mywebsiteconf/www.homeapp.top.key;
		proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header Host $host;
		# additional config
		location / {
            proxy_pass http://$myip:99;
            index  index.html index.htm index.jsp;
			
        }
			error_page 400 401 402 403 405 408 410 412 413 414 415 500 501 502 503 504 506 /505.html;
			location = /502.htm {
				root   html/;
			}
		error_page  404  https://homeapp.top;	 
	}

	# subdomains redirect
	server {
		listen 443 ssl http2;

		server_name *.homeapp.top;
		proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header Host $host;
		# SSL
		ssl_certificate /usr/local/nginx/mywebsiteconf/www.homeapp.top.pem;
		ssl_certificate_key /usr/local/nginx/mywebsiteconf/www.homeapp.top.key;
		
		return 301 https://homeapp.top$request_uri;
		error_page 400 401 402 403 404 405 408 410 412 413 414 415 500 501 502 503 504 506 /505.html;
			location = /502.htm {
				root   html/;
			}
	}

	# HTTP redirect X-Forwarded-For 和 X-Real-IP
	server {
		listen 80;

		server_name .homeapp.top;
		proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header Host $host;
		return 301 https://homeapp.top$request_uri;
	}
```




