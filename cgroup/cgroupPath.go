package cgroup

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

//find hierarchy cgroup root path via "cat /proc/self/mountinfo"
//30 27 0:24 / /sys/fs/cgroup/memory rw, nosuid, nodev, noexec, relatime shared: l3 cgroup cgroup rw, memory
//parm eg. "memory"
func FindCgroupMountpoint(subsystem string) string{
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _,opt := range strings.Split(fields[len(fields) - 1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

//get Subsystem Path
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)
	if _, err := os.Stat (path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err == nil {
			   //do nothing
			}else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	}else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}