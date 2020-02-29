package container

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	ENV_EXEC_CMD = "socker_cmd"
	ENV_EXEC_PID = "socker_pid"
	)

func ExecContainer(containerName string, cmdArray []string) {
	//get pid by name
	pid, err := GetContainerPidByName(containerName)
	if err != nil {
		log.Errorf("Exec container getcontainerPidByName %s error %v", containerName, err)
		return
	}

	cmdStr := strings.Join(cmdArray, " ")
	log.Infof("container pid %s", pid)
	log.Infof("command %s", cmdStr)

	//fork a child process
	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	//attached environment parameters
	os.Setenv(ENV_EXEC_PID, pid)
	os.Setenv(ENV_EXEC_CMD, cmdStr)

	//exec cmd
	//went back to command.go
	if err := cmd.Run(); err != nil {
		log.Errorf("Exec container %s error %v", containerName, err)
	}
}

func GetContainerPidByName(containerName string) (string, error) {
	dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
	configFilePath := dirURL + ConfigName

	//read config
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return "", err
	}

	var containerInfo ContainerInfo
	//json unmarshal
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		return "", err
	}
	return containerInfo.Pid, nil
}