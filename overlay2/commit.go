package overlay2

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func CommitContainer(imageName string){
	mergeDirURL := "/root" + "/" + MERGE
	imageTarURL := "/root" + "/" + imageName + ".tar"
	fmt.Printf("imageTar: %s", imageTarURL)
	cmd := "tar" + " " + "-czf" + " " + imageTarURL + " " + "-C" + " " + mergeDirURL + " " + "."
	fmt.Printf("commit all: %s", cmd)
	if err := exec.Command(cmd); err != nil {

		log.Errorf("Tar folder %s error %v", mergeDirURL, err)
	}
}