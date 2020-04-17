package overlay2

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
)

const (
	WORK       = "/root/workDir/%s"
	MERGE      = "/root/mergeDir/%s"
	//upper layer == volume
	UPPERLAYER = "/root/upperLayer/%s"
	//"/root/busybox"
	LOWER	   = "/root/%s"
	// volume and lower merge
	ROOT 	   = "/root"
)

//mount busybox with overlay2 
//delete function not finish
func NewWorkSpace(containerName string, imageName string) {
	//read-write layer
	CreateUpperLayer(containerName)
	//read-only layer
	CreateLowerLayer(imageName)
	//work layer
	CreateWorkDir(containerName)

	CreateMergeDir(containerName)
	CreateMountPiont(imageName, containerName)
}

func DeleteWorkSpace(containerName string, imageName string){
	//delete mount point
	DeleteMountPoint(containerName)

	DeleteWorkDir(containerName)
	//dont delete upperdir.
	//upperdir play the role of volume in aufs
}

func DeleteMountPoint(containerName string) {
	mergeURL := fmt.Sprintf(MERGE, containerName)
	cmd := exec.Command("umount", "-v", mergeURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		//log.Printf("umount merged cmd all : %s", cmd)
		log.Errorf("umount merged error: %v", err)
	}
	//delete merged directory
	//DeleteMergedDir(rootURL)
}
//merge dir is mnt dir, tar this dir, and get image
//func DeleteMergedDir(rootURL string) {
//	mergedDirURL := rootURL + "/" + MERGE
//	if err := os.RemoveAll(mergedDirURL); err != nil {
//		log.Errorf("Remove dir %s error %v", mergedDirURL, err)
//	}
//}

func DeleteWorkDir(containerName string) {
	workDirURL := fmt.Sprintf(WORK, containerName)
	if err := os.RemoveAll(workDirURL); err != nil {
		log.Errorf("Remove dir %s error %v", workDirURL, err)
	}
}

func CreateMountPiont(imageName, containerName string) {
	dirCmd := "lowerdir=" + LOWER + ",upperdir=" + UPPERLAYER + ",workdir=" + WORK
	dirs := fmt.Sprintf(dirCmd, imageName, containerName, containerName)
	mergeURL := fmt.Sprintf(MERGE, containerName)
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, mergeURL)
	//log.Printf("mount overlay2 command full : %s", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Errorf("Mount overlay error: %v", err)
	}
}

func CreateMergeDir(containerName string) {
	mergeDirURL := fmt.Sprintf(MERGE, containerName)
	bo, err := PathExists(mergeDirURL)
	if err != nil {
		log.Errorf("Find mergeDir error : %v", err)
	}
	if bo == false {
		if err := os.MkdirAll(mergeDirURL, 0777); err != nil {
			log.Errorf("Mkdir dir %s error, %v", mergeDirURL, err)
		}
	}
}

func CreateWorkDir(containerName string) {
	workURL := fmt.Sprintf(WORK, containerName)
	bo, err := PathExists(workURL)
	if err != nil {
		log.Errorf("Find workDir error : %v", err)
	}
	if bo == false {
		if err := os.MkdirAll(workURL, 0777); err != nil {
			log.Errorf("Mkdir dir %s error, %v", workURL, err)
		}
	}
}

func CreateLowerLayer(imageName string) error {
	unTarFolderURL := fmt.Sprintf(LOWER, imageName)
	imageTarURL := ROOT + "/" + imageName + ".tar"
	exist, err := PathExists(unTarFolderURL)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists, %v", unTarFolderURL, err)
		return err
	}
	if exist == false {
		if err := os.MkdirAll(unTarFolderURL, 0777); err != nil {
			log.Errorf("Mkdir dir %s error, %v", unTarFolderURL, err)
			return err
		}

		//zhe li xu yao tiao shi wei shen me bu bao cuo !!!!!!!!!!!!
		if cmd, err := exec.Command("tar", "-xvf", imageTarURL, "-C", unTarFolderURL).CombinedOutput(); err != nil {
			log.Printf("unTar command full : %v", cmd)
			log.Errorf("unTar tar %s error, %v", unTarFolderURL, err)
			return err
		}
	}
	return nil
}

func CreateUpperLayer(containerName string) {
	upperURL := fmt.Sprintf(UPPERLAYER, containerName)
	bo, err := PathExists(upperURL)
	if err != nil {
		log.Errorf("Find upperDir error : %v", err)
	}
	if bo == false {
		if err := os.MkdirAll(upperURL, 0777); err != nil {
			log.Errorf("Mkdir dir %s error, %v", upperURL, err)
		}
	}
}

func PathExists(url string) (bool, error) {
	_, err := os.Stat(url)
	if err == nil {
		return true, nil
	}
	//attention to the parameter
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}


