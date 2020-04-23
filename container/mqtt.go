package container

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
)

func MqttClient(mqtt bool, containerName string) {
	//call mqtt watcher
	if mqtt {
		cmd := exec.Command("/bin/sh", "-c", "sudo -b nohup ./mqttWatcher start "+containerName)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Errorf("call mqtt watcher error %v", err)
		}
	}
}


func StopMqtt(containerName string) {
	cmd := exec.Command("/bin/sh", "-c",	"sudo ./mqttWatcher stop " + containerName)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("stop mqtt watcher error %v", err)
	}
}