package container

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

//create a parent process for container command
//return that command (it needs Start() function to run)
func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	//run this process itself with args
	cmd := exec.Command("/proc/self/exe",args...)
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
	return cmd
}

//called by InitCommand
//the first process runs by new container
//mount proc file-system so you can check process source with "ps"
func InitProcess(command string, args []string) error{
	log.Infof("command %s", command)
	log.Print("start to mount")

	// systemd 加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示
	//声明你要这个新的mount namespace独立。
	syscall.Mount("", "/", "", syscall.MS_PRIVATE | syscall.MS_REC, "")
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	//exec会执行参数指定的命令，但是并不创建新的进程，只在当前进程空间内执行，即替换当前进程的执行内容，他们重用同一个进程号PID。
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}