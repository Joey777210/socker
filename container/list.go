package container

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"text/tabwriter"
)

func ListContainers(){
	dirURL := fmt.Sprintf(DefaultInfoLocation, "")
	dirURL = dirURL[:len(dirURL)-1]
	//read all file
	files, err := ioutil.ReadDir(dirURL)
	if err != nil {
		log.Errorf("Read dir %s error %v", dirURL, err)
		return
	}

	var containers []*ContainerInfo
	//range all files in this dir
	for _, file := range files {

		tmpContainer, err := getContainerInfo(file)
		if err != nil {
			log.Errorf("Get container info error %v", err)
			continue
		}
		containers = append(containers, tmpContainer)
	}

	//tabwriter.NewWriter print contianer information on console
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	//console info column
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, item := range containers {
		fmt.Fprintf(w,"%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id,
			item.Name,
			item.Pid,
			item.Status,
			item.Command,
			item.CreatedTime)
	}

	//flush stdout stream buffer, print the container list
	if err := w.Flush(); err != nil {
		log.Errorf("Flush error %v", err)
		return
	}
}

func getContainerInfo(file os.FileInfo) (*ContainerInfo, error) {
	//get file name
	containerName := file.Name()
	//generate file absolute path base on file name
	configFileDir := fmt.Sprintf(DefaultInfoLocation, containerName)
	configFileDir = configFileDir + ConfigName
	//read config.json
	content, err := ioutil.ReadFile(configFileDir)
	if err != nil {
		log.Errorf("Read file %s error %v", configFileDir, err)
		return nil, err
	}

	var containerInfo ContainerInfo
	//json unmarshal to containerInfo class
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		log.Errorf("Json unmarshal err %v", err)
		return nil, err
	}
	return &containerInfo, nil
}
