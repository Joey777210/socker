package container

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func LogContainer(containerName string) {
	//find dir
	dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
	LogFileLocation := dirURL + "container.log"
	//open log
	file, err := os.Open(LogFileLocation)
	defer file.Close()
	if err != nil {
		log.Errorf("Log container open file %s error %v", LogFileLocation, err)
		return
	}
	//read log file
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorf("Log container read file %s error %v", LogFileLocation, err)

	}
	//output content
	fmt.Fprint(os.Stdout, string(content))
}
