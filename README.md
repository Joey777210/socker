# Socker  
## Socker是什么
  Socker是仿Docker的Linux容器引擎，包括容器运行、管理、网络连接等部分。基于阿里大佬xianlubird的书《自己动手写Docker》和runc的源代码实现。
  过程中参照很多大佬博客，没有详细记录，统一远程感谢。
## Get Start
*环境*  
  `Ubuntu 18.04，go 10.4`  
**下载软件*  
  ```
	git clone https://github.com/Joey777210/Socker.git
	go get
	go build
  ```
*运行*  
从Docker中拷贝出打包出一个ubuntu.tar的image，放在`/root`目录下  
```
	sudo ./Socker run -ti --name ubuntu ubuntu sh
```

## 使用指南  
### Run Socker    
```
./Socker run -ti IMAGENAME COMMAND  
如： sudo ./Socker run -ti ubuntu sh
```
### 后台运行  
```
 run -d
e.g. sudo ./Socker run -d ubuntu top -b  
```
### 命名  
```
--name NAME
e.g. sudo ./Socker run -ti --name myContainer ubuntu sh    
```
### 指定环境变量运行容器  
```
-e env  
e.g.   sudo ./Socker run -ti --name socker -e bird=123 -e luck=bird ubuntu sh  
```

### 查看正在运行的容器    
`./Socker ps`  
### 查看容器日志  
`./Socker logs ContainerName`  
### 进入后台运行的容器  
`./Socker exec NAME sh`  
### 停止一个容器  
`./Socker stop NAME`  
### 删除一个容器  
`./Socker rm NAME`  

### 容器资源管理  
1.memory  
-----------
> `-m 100m`  
2.cpushare
 ------------- 
> `-cpushare 512`  
3.cpuset  
 -----------
> `-cpuset 1`  
### *通过容器制作镜像*
>在一个Terminal上运行容器
>打开另一个Terminal并运行命令
`./Socker commit IMAGENAME`  
>现在你可以看到 `/root` 目录下生成了镜像文件 `IMAGENAME.tar`  

### *查看所有镜像*  
`sudo ./Socker image -ls`  
### *删除镜像*  
`sudo ./Socker image -rm IMAGENAME`  

## 网络
### *创建一个网桥*
`sudo ./Socker network create --driver bridge --subnet 192.168.10.1/24 BridgeName`  
### *利用网桥使容器能够接入互联网*
`sudo ./Socker run -ti -net BridgeName ubuntu sh`  
现在可以使用 `ping` 命令测试你的容器了  
### *列出所创建的网络*
`sudo ./Socker network list`    

## 使用  
1. 下载Socker和SockerMQTTWatcher后，放在$GOPATH/src下   
2. 使用docker下载ubuntu.tar 放在go/src下
3. 进入Socker目录中，make get
4. make build

## 目前遇到的问题和bug：  
  1.在busybox中mount /proc时会报错。mount2 point  
