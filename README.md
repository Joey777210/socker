# Socker  
## Socker是什么
Socker是仿Docker的Linux容器引擎，包括容器运行、管理、网络连接等部分。基于阿里大佬xianlubird的书《自己动手写Docker》和runc的源代码实现。  

过程中参照很多大佬博客，没有详细记录，统一远程感谢。
## Get Start
### 环境
>Ubuntu 18.04.4 LTS 内核版本5.3.0-51-generic
>Go version go1.13.8 Linux/amd64
>C gcc version 7.5.0

### 下载软件 
```
git clone https://github.com/Joey777210/Socker.git
```
### 编译
>将socker和SockerMQTTWatcher放在$GOPATH/src下
>进入socker目录下，使用make get命令，下载所需的Go依赖库
>使用make build命令，编译项目
>现在你可以使用socker了！

### 运行  

```
sudo ./Socker run -ti --name myubuntu ubuntu sh
```


## 使用  
### 注意：您需要在root权限或sudoer中利用sudo命令使用socker  
* ### run –ti/d [镜像名称] [启动命令]： 运行容器  
参数：-ti				前台运行容器  
```
    如：	sudo socker run -ti ubuntu sh  
```
参数：-d				后台运行容器  
```
 		如：	sudo socker run -d ubuntu top -b    
```
参数：--name			容器命名
```
		如：	sudo socker run -d --name myubuntu Ubuntu top -b   
```
参数：-e				指定环境变量运行容器
```
    如：	sudo socker run –ti -e bird=123 -e luck=bird ubuntu sh    
```
参数：-m				限制容器内存资源  
```
	  如：	sudo socker run -ti –m 100m ubuntu sh  
```
参数：-cpuset			限制容器使用CPU核心数  
```
  	如：	sudo socker run –ti –cpuset 1 ubuntu sh  
```
参数：-cpushare		限制容器使用CPU时间片的权重  
```
  	如：	sudo socker run –ti –cpushare 512 ubuntu sh  
```
#### 创建网络后可用的参数：  
参数：-net				配置容器网络  
```
  	如：	sudo socker run –ti –net testbridge ubuntu sh  
```
#### 注：此时还没有配置DNS，无法按域名访问，可以使用访问外网IP地址  
参数：-p				配置容器端口映射  
```
  	如：	sudo socker run –ti –net testbridge –p 8080：80 ubuntu sh  
```
* ### ps：	查看容器列表   
```
  	如：	sudo socker ps  
```
* ### logs [容器名称]:	查看容器日志  
```
  	如：	sudo socker logs myubuntu  
```
* ### exec [容器名称] [运行命令]:	进入后台运行容器的Namespace  
```
   	如：	sudo socker exec myubuntu sh      
```
* ### stop [容器名称]	:	停止指定容器  
```
    如：	sudo socker stop NAME  
```
* ### rm [容器名称]:		删除指定容器（已停止的）  
```
 		如：	sudo socker rm myubuntu    
```
* ### commit [容器名称]：制作指定容器的镜像  
```
  	如：	sudo socker commit myubuntu	   
```
* ### image：	管理镜像  
	参数：-ls				查看镜像列表  
```
 		如：	sudo socker image -ls   
```
  参数：-rm				删除指定镜像   
```
 		如：	sudo socker image -rm myubuntu   
```
* ### network：	管理网络  
	参数：create			创建网络  
    二级参数：--driver		配置子网络驱动（网桥）  
		    	   --subnet		配置子网络网段和子网掩码  
```
  	如：	 sudo socker network create --driver bridge --subnet 192.168.10.1/24 testbridge  
```
参数：list				列出所创建的网络  
```
 		如：	sudo socker network list   
```
## 帮助  
参数：-h/help			查看帮助  
```
  	如：	sudo socker–h  
		      sudo SockerMQTTWatcher–h  
```
