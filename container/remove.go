package container

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
)

func (c *Container) RemoveContainer(containerName string) {
	containerInfo, err := GetContainerInfoByName(containerName)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerName, err)
		return
	}

	if containerInfo.Status != STOP {
		log.Errorf("Can't remove running container")
		return
	}

	//find file path
	dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove container file %s error %v", containerName, err)
		return
	}
}
