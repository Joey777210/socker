package container

import (
	"time"
)

type Container struct {
	Pid         string `json:"pid"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	CreatedTime string `json:"createTime"`
	Status      string `json:"status"`
	//Volume      string `json:"volume"`     //容器的数据卷, store upper layer
	PortMapping []string `json:"portmapping"` //端口映射
}

func NewContainer(containerName string) *Container {
	//随机数生成containerID
	id := randStringBytes(10)
	createTime := time.Now().Format("2006-01-02 15:04:05")

	if containerName == "" {
		containerName = id
	}

	return &Container{
		Id:          id,
		Name:        containerName,
		CreatedTime: createTime,
	}
}
