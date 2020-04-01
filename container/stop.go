package container

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"strconv"
	"syscall"
)

func StopContainer(containerName string) {
	//get PID
	pid, err := GetContainerPidByName(containerName)
	if err != nil {
		log.Errorf("Get container pid by name %s error %v", containerName, err)
		return
	}
	//PID in string to int
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		log.Errorf("Conver pid from string to int error %v", err)
		return
	}
	//system call kill send signal to process
	if err := syscall.Kill(pidInt, syscall.SIGTERM); err!= nil {
		log.Errorf("Stop container %s error %v", containerName, err)
		return
	}
	//get info by name
	containerInfo, err := GetContainerInfoByName(containerName)
	if err != nil {
		log.Errorf("get container %s info error %v", containerName, err)
		return
	}
	//change container status
	containerInfo.Status = STOP
	containerInfo.Pid = " "
	//json marshal
	newContentBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Json marshal %s error %v", containerName, err)
		return
	}
	dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
	configFilePath := dirURL + ConfigName
	if err := ioutil.WriteFile(configFilePath, newContentBytes, 0622); err != nil {
		log.Errorf("Write file %s error", configFilePath, err)
	}
}


func GetContainerInfoByName(containerName string) (*ContainerInfo, error) {
	dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
	configFilePath := dirURL + ConfigName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Errorf("Read file %s error %v", configFilePath, err)
		return nil, err
	}

	var containerInfo ContainerInfo
	//json unmarshal
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		log.Errorf("getContainerInfoByName unmarshal error %v", err)
		return nil, err
	}
	return &containerInfo, nil
}