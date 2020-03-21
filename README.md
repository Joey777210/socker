# Socker  
## What Socker is?  
  Socker is container engine for linux, I wrote this project for study how Docker works. Socker is based on 《自己动手写Docker》 written by xianlubird and runc.
## How to start?
×System environment: 
  `Ubuntu 18.04，go 10.4`
*Download the source code   
  ```
	go get
	go build
	./Socker run -ti sh
  ```
If you can see bash running, it means you got it right.

## How to use Socker?  
### *Run Socker   
```
./Socker run -ti sh
```
### *Run container background  
`./Socker run -d command`  
### *Name a container  
`./Socker run -ti --name NAME`  
### *Check running container  
`./Socker ps`  
### *Ckeck container logs  
`./Socker logs ContainerName`  
### *Enter a container    
`./Socker exec NAME sh`  
### *Stop a container   
`./Socker stop NAME`  
### *Remove a container   
`./Socker rm NAME`  

### *Container resource constrains
1.memory  
`./Socker run -ti -m 100m`  
2.cpushare  
`./Socker run -ti -cpushare 512`  
3.cpuset   
`./Socker run -ti -cpuset 1`  
### *Image packaging  
>Run container in one Terminal  
>Open another Terminal and use command  
`./Socker commit image`  
>Then you can get a commited image under /root  

## Network  
### *Create a bridge network  
`./Socker network create --driver bridge --subnet 192.168.10.1/24 BridgeName`  
### *Use bridge  connect to Internet  
`./Socker run -ti -net BridgeName sh`  
Now your container is able to `ping`   


## 目前思路：结合mqtt实现容器内的数据下发  
## 目前遇到的问题和bug：  
  1.在busybox中mount /proc时会报错。mount2 point  
  2.后台运行top时， 通过ps -ef并不能看到init接管top进程  
