package command

import (
	"Socker/cgroup"
	"Socker/container"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

//called by runCommand
func Run(tty bool, command []string, resourceConfig cgroup.ResourceConfig, containerName string){
	//gets the command
	parent, writePipe:= container.NewParentProcess(tty, containerName)

	if parent == nil {
		log.Errorf("New parent process error")
		return
	}

	//parent command starts to manipulate
	if err := parent.Start(); err != nil{
		log.Error(err)
	}

	containerName, err := container.RecordContainerInfo(parent.Process.Pid, command, containerName)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return
	}

	//use socker-cgroup as cgroup name
	//create cgroup manager and set res
	cgroupManager := cgroup.NewCgroupManager("socker-cgroup")
	defer cgroupManager.Destroy()
	//set res
	cgroupManager.Set(&resourceConfig)
	//set container into cgroup
	cgroupManager.Apply(parent.Process.Pid)
	//init container
	sendInitCommand(command, writePipe)
	log.Print("socker: exit socker")
	if tty {
		parent.Wait()
		container.DeleteContainerInfo(containerName)
	}

	//create image related
	//overlay2.DeleteWorkSpace("/root", "/root/mergeDir")

	//os.Exit(-1)
}

func sendInitCommand(command []string, writePipe *os.File) {
	cmdStr := strings.Join(command, " ")
	log.Infof("command init writePipe is %s", cmdStr)
	writePipe.WriteString(cmdStr)
	writePipe.Close()
}


