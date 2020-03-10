package container

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type ContainerInfo struct {
	Pid			string `json:"pid"`
	Id			string `json:"id"`
	Name		string `json:"name"`
	Command		string `json:"command"`
	CreatedTime	string `json:"createTime"`
	Status		string `json:"status"`
	Volume      string `json:"volume"`     //容器的数据卷
	PortMapping []string `json:"portmapping"` //端口映射
}

var (
	RUNNING				= "running"
	STOP				= "stopped"
	Exit				= "exited"
	DefaultInfoLocation	= "/var/run/socker/%s/"
	ConfigName			= "config.json"
)

func RecordContainerInfo(containerPID int, commandArray []string, containerName string, containerID string) (string, error){
	//1, generate num container ID
	id := randStringBytes(10)

	createTime := time.Now().Format("2006-01-02 21:01:05")
	command := strings.Join(commandArray, " ")

	if containerName == "" {
		containerName = id
	}

	containerInfo := &ContainerInfo{
		Pid:         strconv.Itoa(containerPID),
		Id:          id,
		Name:        containerName,
		Command:     command,
		CreatedTime: createTime,
		Status:      RUNNING,
	}

	//json to string
	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirURL, 0622); err != nil {
		log.Errorf("Mkdir error %s error %v", dirURL, err)
		return "", err
	}

	fileName := dirURL + "/" + ConfigName

	//create config.json
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Errorf("Create file %s error %v", fileName, err)
		return "", err
	}
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("File write string error %v", err)
		return "", err
	}

	return containerName, nil
}


//generate a container ID
func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func DeleteContainerInfo(containerId string){
	dirURL := fmt.Sprintf(DefaultInfoLocation, containerId)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove dir %s error %v", dirURL, err)
	}
}