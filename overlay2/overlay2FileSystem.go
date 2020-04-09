package overlay2

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
)

const (
	BUSYBOX    = "ubuntu"
	BUSYBOXTAR = "ubuntu.tar"
	WORK       = "workDir"
	MERGE      = "mergeDir"
	UPPERLAYER = "upperLayer"
)

//mount busybox with overlay2 
//delete function not finish
func NewWorkSpace(rootURL string, mergedURL string) {
	//read-write layer
	CreateUpperLayer(rootURL)
	//read-only layer
	CreateLowerLayer(rootURL)
	//read-write layer
	CreateWorkDir(rootURL)

	CreateMergeDir(rootURL)
	CreateMountPiont(rootURL, mergedURL)
}

func DeleteWorkSpace(rootURL string, mergedURL string){
	//delete mount point
	DeleteMountPoint(rootURL, mergedURL)

	DeleteWorkDir(rootURL)
	//dont delete upperdir.
	//upperdir play the role of volume in aufs
}

func DeleteMountPoint(rootURL string, mergedURL string) {
	cmd := exec.Command("umount", "-v", mergedURL)
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

func DeleteWorkDir(rootURL string) {
	workDirURL := rootURL + "/" + WORK
	if err := os.RemoveAll(workDirURL); err != nil {
		log.Errorf("Remove dir %s error %v", workDirURL, err)
	}
}

func CreateMountPiont(rootURL string, mergedURL string) {
	dirs := "lowerdir=" + rootURL + BUSYBOX + ",upperdir=" + rootURL + UPPERLAYER + ",workdir=" + rootURL + WORK
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, rootURL + MERGE)
	//log.Printf("mount overlay2 command full : %s", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Errorf("Mount overlay error: %v", err)
	}
}

func CreateMergeDir(rootURL string) {
	mergeDirURL := rootURL + MERGE
	bo, err := PathExists(mergeDirURL)
	if err != nil {
		log.Errorf("Find mergeDir error : %v", err)
	}
	if bo == false {
		if err := os.Mkdir(mergeDirURL, 0777); err != nil {
			log.Errorf("Mkdir dir %s error, %v", mergeDirURL, err)
		}
	}
}

func CreateWorkDir(rootURL string) {
	workDirURL := rootURL + WORK
	bo, err := PathExists(workDirURL)
	if err != nil {
		log.Errorf("Find workDir error : %v", err)
	}
	if bo == false {
		if err := os.Mkdir(workDirURL, 0777); err != nil {
			log.Errorf("Mkdir dir %s error, %v", workDirURL, err)
		}
	}
}

func CreateLowerLayer(rootURL string) {
	busyboxURL := rootURL + BUSYBOX
	busyboxTarURL := rootURL + BUSYBOXTAR
	exist, err := PathExists(busyboxURL)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists, %v", busyboxURL, err)
	}
	if exist == false {
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			log.Errorf("Mkdir dir %s error, %v", busyboxURL, err)
		}

		//zhe li xu yao tiao shi wei shen me bu bao cuo !!!!!!!!!!!!
		if cmd, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			log.Printf("unTar command full : %v", cmd)
			log.Errorf("unTar tar %s error, %v", busyboxTarURL, err)
		}
	}
}

func CreateUpperLayer(rootURL string) {
	upperURL := rootURL + UPPERLAYER
	bo, err := PathExists(upperURL)
	if err != nil {
		log.Errorf("Find upperDir error : %v", err)
	}
	if bo == false {
		if err := os.Mkdir(upperURL, 0777); err != nil {
			log.Errorf("Mkdir dir %s error, %v", upperURL, err)
		}
	}
}

func PathExists(busyboxURL string) (bool, error) {
	_, err := os.Stat(busyboxURL)
	if err == nil {
		return true, nil
	}
	//attention to the parameter
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}


