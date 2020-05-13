package container

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
)

func (c *Container) RemoveContainer() {
	containerInfo, err := GetContainerInfoByName(c.Name)
	if err != nil {
		log.Errorf("Get container %s info error %v", c.Name, err)
		return
	}

	if containerInfo.Status != STOP {
		log.Errorf("Can't remove running container")
		return
	}

	//find file path
	dirURL := fmt.Sprintf(DefaultInfoLocation, c.Name)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove container file %s error %v", c.Name, err)
		return
	}
}
