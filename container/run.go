package container

import (
	"Socker/cgroup"
	"Socker/network"
	"Socker/overlay2"
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"os"
	"strings"
	"time"
)

func (c *Container) Run(tty bool, command []string, resourceConfig *cgroup.ResourceConfig, nw string, mqtt bool, imageName string, envSlice []string, portMapping []string) {
	//gets the command
	parent, writePipe := c.NewParentProcess(tty, imageName, envSlice)

	if parent == nil {
		log.Errorf("New parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	err := c.RecordContainerInfo(command, parent.Process.Pid, portMapping)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return
	}

	//use socker-cgroup as cgroup name
	//create cgroup manager and set res
	cgroupManager := cgroup.NewCgroupManager(DefaultLocation + "/socker-cgroup")
	defer cgroupManager.Destroy()
	//set res
	cgroupManager.Set(resourceConfig)
	//set container into cgroup
	cgroupManager.Apply(parent.Process.Pid)

	if nw != "" {
		// 配置网络连接
		network.Init()
		containerInfo := &ContainerInfo{
			Pid:         c.Pid,
			Id:          c.Id,

			PortMapping: c.PortMapping,
		}
		if err := network.Connect(nw, containerInfo); err != nil {
			log.Errorf("Error Connect Network %v", err)
			return
		}
	}

	if mqtt {
		mqttManager := NewMqttManager(c.Name)
		go mqttManager.Create()
	}

	//init container
	sendInitCommand(command, writePipe)
	if tty {
		parent.Wait()
		if mqtt {
			mqttManager := NewMqttManager(c.Name)
			mqttManager.Stop()
		}
		DeleteContainerInfo(c.Name)
	}
	//create image related
	overlay2.DeleteWorkSpace(c.Name)
	os.Exit(0)
}

func sendInitCommand(command []string, writePipe *os.File) {
	cmdStr := strings.Join(command, " ")
	log.Infof("command init writePipe is %s", cmdStr)
	writePipe.WriteString(cmdStr)
	writePipe.Close()
}

//生成随机的containerID
func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
