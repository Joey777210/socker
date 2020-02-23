package container

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

//create a parent process for container
//return that command (it needs Start() function to run)
func NewParentProcess(tty bool) (*exec.Cmd, *os.File){
	readPipe, writePipe, err := NewPipe()
	if err != nil{
		log.Errorf("new pipe error %v", err)
		return nil, nil
	}

	//run this process itself with args
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:syscall.CLONE_NEWIPC | syscall.CLONE_NEWUSER | syscall.CLONE_NEWPID |
		syscall.CLONE_NEWNET | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS,
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
	}
	cmd.ExtraFiles = []*os.File{readPipe};
	return cmd, writePipe
}

//called by InitCommand
//the first process runs by new container
//mount proc file-system so you can check process source with "ps"
func InitProcess() error{
	cmdArray := readUserCommand()
	log.Infof("command in init is %v", cmdArray)

	// systemd 加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示
	//声明你要这个新的mount namespace独立。
	syscall.Mount("", "/", "", syscall.MS_PRIVATE | syscall.MS_REC, "")
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := cmdArray
	//find absolute path of command
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("find command error %v", err)
		return err
	}
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