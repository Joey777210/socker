package cgroup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CpusetSubSystem struct {

}

//return the name of current subsystem
func (s *CpusetSubSystem)Name() string{
	return "cpuset"
}

//set res to cgroup
func (s *CpusetSubSystem)Set (cgroupPath string, res *ResourceConfig) error {
	if SubsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil{
		if res.CpuSet != ""{
			if err := ioutil.WriteFile(path.Join(SubsysCgroupPath, "cpuset.cpus"), []byte(res.CpuSet), 0644); err != nil {
				return fmt.Errorf("set cgroup cpuset fail %v", err)
			}
		}
		return nil
	}else {
		return err
	}
}

//apply cgroup to pid process || add pid to cgroup
func (s *CpusetSubSystem)Apply(cgroupPath string, pid int) error {
	if SubsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if err := ioutil.WriteFile(path.Join(SubsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	}else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}

//remove a cgroup in parmPath
func (s *CpusetSubSystem)Remove(cgroupPath string) error {
	if SubsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil{
		return os.Remove(SubsysCgroupPath)
	}else {
		return err
	}

}
