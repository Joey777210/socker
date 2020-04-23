package command

import (
	"Socker/cgroup"
	"Socker/container"
	"Socker/network"
	"Socker/overlay2"
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

//called by runCommand
func Run(tty bool, command []string, resourceConfig *cgroup.ResourceConfig, containerName string, nw string, portmapping []string, mqtt bool, imageName string, envSlice []string){

	containerID := randStringBytes(10)
	//gets the command
	parent, writePipe:= container.NewParentProcess(tty, containerName, imageName, envSlice)

	if parent == nil {
		log.Errorf("New parent process error")
		return
	}

	//parent command starts to manipulate
	if err := parent.Start(); err != nil{
		log.Error(err)
	}

	containerName, err := container.RecordContainerInfo(parent.Process.Pid, command, containerName, containerID)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return
	}

	//use socker-cgroup as cgroup name
	//create cgroup manager and set res
	cgroupManager := cgroup.NewCgroupManager("socker-cgroup")
	defer cgroupManager.Destroy()
	//set res
	cgroupManager.Set(resourceConfig)
	//set container into cgroup
	cgroupManager.Apply(parent.Process.Pid)

	if nw != "" {
		// config container network
		network.Init()
		containerInfo := &container.ContainerInfo{
			Id:          containerID,
			Pid:         strconv.Itoa(parent.Process.Pid),
			Name:        containerName,
			PortMapping: portmapping,
		}
		if err := network.Connect(nw, containerInfo); err != nil {
			log.Errorf("Error Connect Network %v", err)
			return
		}
	}


	go container.MqttClient(mqtt, containerName)

	//init container
	sendInitCommand(command, writePipe)
	if tty {
		parent.Wait()
		if mqtt {
			container.StopMqtt(containerName)
		}
		container.DeleteContainerInfo(containerName)
	}
	//create image related
	overlay2.DeleteWorkSpace(containerName, imageName)
	os.Exit(0)
}

func sendInitCommand(command []string, writePipe *os.File) {
	cmdStr := strings.Join(command, " ")
	log.Infof("command init writePipe is %s", cmdStr)
	writePipe.WriteString(cmdStr)
	writePipe.Close()
}


func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
