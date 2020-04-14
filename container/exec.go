package container

import (
	_ "Socker/nsenter"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
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
	cmd.Env = append(os.Environ(), getEnvsByPid(pid)...)
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

func getEnvsByPid(pid string) []string {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Read file %s error %v", path, err)
		return nil
	}
	//env split by \u0000
	envs := strings.Split(string(contentBytes), "\u0000")
	return envs
}