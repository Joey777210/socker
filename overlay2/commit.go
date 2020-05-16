package overlay2

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os/exec"
)

func CommitContainer(imageName string){
	mergeDirURL := fmt.Sprintf(MERGE, imageName)
	imageTarURL := "/root" + "/" + imageName + ".tar"
	fmt.Printf("imageTar: %s", imageTarURL)
	//cmd := "tar" + " " + "-czPf" + " " + imageTarURL + " " + mergeDirURL
	//cmd := "/bin/sh " + "-c " + "sudo " + "tar" + " " + "-czPf" + " " + imageTarURL + " " + mergeDirURL
	//cmd := exec.Command("/bin/sh", "-c", "sudo " + "tar" + " " + "-czPf" + " " + imageTarURL + " " + mergeDirURL)
	//fmt.Printf("commit all: %s", cmd)
	if _, err := exec.Command("tar", "-czf", imageTarURL, mergeDirURL).CombinedOutput(); err != nil {
		log.Errorf("Tar folder %s error %v", mergeDirURL, err)
	}
}