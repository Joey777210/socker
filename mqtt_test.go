package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"testing"
)

func TestMqtt(t *testing.T) {
	containerName := "sockerss"
	str := "sudo -b nohup ../mqttWatcher start "+ containerName
	cmd := exec.Command(str)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("call mqtt watcher error %v", err)
	}
}
