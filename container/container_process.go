package container

import (
	"Socker/overlay2"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const (
	WORKDIR = "/root/mergeDir"
	ROOT = "/root/"
)


//create a parent process for container
//return that command (it needs Start() function to run)
func NewParentProcess(tty bool, containerName string) (*exec.Cmd, *os.File){
	readPipe, writePipe, err := NewPipe()
	if err != nil{
		log.Errorf("new pipe error %v", err)
		return nil, nil
	}
	//run this process itself with args
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:unix.CLONE_NEWIPC | unix.CLONE_NEWPID |
			unix.CLONE_NEWNET | unix.CLONE_NEWUTS | unix.CLONE_NEWNS | unix.CLONE_NEWUSER,
		//found in github issue. solve mount /proc problem
		UidMappings:[]syscall.SysProcIDMap{
			{ ContainerID: 0, HostID: 0, Size: 1, },
		},
		GidMappings:[]syscall.SysProcIDMap{
			{ ContainerID: 0, HostID: 0, Size: 1, },
		},
	}
	//cmd.SysProcAttr.Credential = &syscall.Credential{Uid:uint32(1), Gid:uint32(1)}

	//tty means whether its a interactive process.
	if tty{
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

	}else {
		//out put log into container.log
		dirURL := fmt.Sprintf(DefaultInfoLocation, containerName)
		if err := os.MkdirAll(dirURL, 0622); err != nil {
			log.Errorf("NewParentProcess mkdir %s error %v", dirURL, err)
			return nil, nil
		}
		stdLogFilePath := dirURL + "container.log"
		stdLogFile, err := os.Create(stdLogFilePath)
		if err != nil {
			log.Errorf("NewParentProcess create file %s error %v", stdLogFilePath, err)
			return nil, nil
		}
		cmd.Stdout = stdLogFile
	}
	cmd.ExtraFiles = []*os.File{readPipe}

	//create image...
	
	overlay2.NewWorkSpace(ROOT, WORKDIR)
	cmd.Dir = WORKDIR
	//cmd.Dir = "/home/joey/go/busybox"

	return cmd, writePipe
}

//called by InitCommand
//the first process runs by new container
//mount proc file-system so you can check process source with "ps"
func InitProcess() error{
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error, cmdArray is nil")
	}
	log.Infof("command in init is %v", cmdArray)

	//new rootfs
	SetUpMount()

	//find absolute path of command
	log.Infof("find command is %s", cmdArray[0])
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("find command error %v", err)
		return err
	}
	argv := cmdArray[0:]
	log.Printf("Found path %s\n", path)

	//exec会执行参数指定的命令，但是并不创建新的进程，只在当前进程空间内执行，即替换当前进程的执行内容，他们重用同一个进程号PID。
	if err := syscall.Exec(path, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

//read cmdStr from read pipe and split with " "
func readUserCommand() []string {
	readPipe := os.NewFile(uintptr(3), "pipe")
	defer readPipe.Close()
	cmdBytes, err := ioutil.ReadAll(readPipe)

	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	cmdStr := string(cmdBytes)
	log.Printf("read pipe get param is %s", cmdStr)
	return strings.Split(cmdStr, " ")

}

//create a pipe between user command and process
func NewPipe()(*os.File, *os.File, error){
	read, write, err := os.Pipe()
	if err != nil{
		return nil, nil, err
	}
	return read, write, nil
}