package container

import (
	"socker/cgroup"
	"socker/network"
	"socker/overlay2"
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"os"
	"strings"
	"time"
)

func (c *Container) Run(tty bool, command []string, resourceConfig *cgroup.ResourceConfig, nw string, imageName string, envSlice []string, portMapping []string) {
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
		network.GetContainerInfo(c.Id, c.Pid, c.PortMapping)
		if err := network.Connect(nw); err != nil {
			log.Errorf("Error Connect Network %v", err)
			return
		}
	}

	//init container
	sendInitCommand(command, writePipe)
	if tty {
		parent.Wait()
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
