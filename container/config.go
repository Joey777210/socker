package container

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"strconv"
)

type ContainerInfo struct {
	Pid         string   `json:"pid"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Command     string   `json:"command"`
	CreatedTime string   `json:"createTime"`
	Status      string   `json:"status"`
	Volume      string   `json:"volume"`      //容器的数据卷, store upper layer
	PortMapping []string `json:"portmapping"` //端口映射
}

var (
	RUNNING             = "running"
	STOP                = "stopped"
	Exit                = "exited"
	DefaultLocation     = "/var/run/socker"
	DefaultInfoLocation = "/var/run/socker/%s/"
	ConfigName          = "config.json"
)

func (c *Container) RecordContainerInfo(containerPID int) error {

	c.Pid = strconv.Itoa(containerPID)
	c.Status = RUNNING

	//json to string
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return err
	}
	jsonStr := string(jsonBytes)

	dirURL := fmt.Sprintf(DefaultInfoLocation, c.Name)
	if err := os.MkdirAll(dirURL, 0622); err != nil {
		log.Errorf("Mkdir error %s error %v", dirURL, err)
		return err
	}

	fileName := dirURL + "/" + ConfigName

	//create config.json
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Errorf("Create file %s error %v", fileName, err)
		return err
	}
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("File write string error %v", err)
		return err
	}

	return nil
}

func DeleteContainerInfo(containerId string) {
	dirURL := fmt.Sprintf(DefaultInfoLocation, containerId)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove dir %s error %v", dirURL, err)
	}
}
