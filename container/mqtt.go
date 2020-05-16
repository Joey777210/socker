package container

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
)

type MqttManager struct {
	containerName string
}

func NewMqttManager(containerName string) *MqttManager {
	return &MqttManager{
		containerName: containerName,
	}
}

func (m *MqttManager) Create() {

	//调用mqttWatcher
	cmd := exec.Command("/bin/sh", "-c", "sudo -b nohup ./mqttWatcher start "+ m.containerName)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("call mqtt watcher error %v", err)
	}
}

func (m *MqttManager) Stop() {
	cmd := exec.Command("/bin/sh", "-c", "sudo ./mqttWatcher stop "+ m.containerName)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("stop mqtt watcher error %v", err)
	}
}
