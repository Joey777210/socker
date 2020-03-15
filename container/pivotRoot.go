package container

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"syscall"
)

func pivotRoot(root string) error {
	if err := syscall.Mount(root, root , "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}

	//create pivotDir to store old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil{
		return nil
	}

	//pivot to new rootfs
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}

	//change to root "/" dir
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}


	pivotDir = filepath.Join("/", ".pivot_root")
	//umount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}

	//delete temporary dir
	return os.Remove(pivotDir)
}


/**
	init mount rootfs
 */
func SetUpMount(){
	//get current path
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Get current location error %v", err)
		return
	}
	log.Infof("Current location is %s", pwd)
	//error: when use this function, you cant mount /proc in busybox...
	//dont know how to solve
	if err = pivotRoot(pwd); err != nil {
		log.Errorf("pivot root error: %v", err)
	}

	// systemd 加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示
	//声明你要这个新的mount namespace独立。
	syscall.Mount("", "/", "", syscall.MS_PRIVATE | syscall.MS_REC, "")

	//mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err2 := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err2 != nil {
		log.Errorf("mount2 error : %v", err2)
	}
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID | syscall.MS_STRICTATIME, "mode=755")
}
