package container

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"strconv"
	"syscall"
)

func (c *Container) StopContainer() {
	//get PID
	pid, err := GetContainerPidByName(c.Name)
	if err != nil {
		log.Errorf("Get container pid by name %s error %v", c.Name, err)
		return
	}
	//PID in string to int
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		log.Errorf("Conver pid from string to int error %v", err)
		return
	}
	//system call kill send signal to process
	if err := syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		log.Errorf("Stop container %s error %v", c.Name, err)
		return
	}
	//get info by name
	containerInfo, err := GetContainerInfoByName(c.Name)
	if err != nil {
		log.Errorf("get container %s info error %v", c.Name, err)
		return
	}
	//change container status
	containerInfo.Status = STOP
	containerInfo.Pid = " "
	//json marshal
	newContentBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Json marshal %s error %v", c.Name, err)
		return
	}
	dirURL := fmt.Sprintf(DefaultInfoLocation, c.Name)
	configFilePath := dirURL + ConfigName
	if err := ioutil.WriteFile(configFilePath, newContentBytes, 0622); err != nil {
		log.Errorf("Write file %s error", configFilePath, err)
	}

	//停止mqttWatcher
	mqttManager := NewMqttManager(c.Name)
	mqttManager.Stop()
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
