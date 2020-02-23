package cgroup

import (
	log "github.com/sirupsen/logrus"
)

type CgroupManager struct {
	//cgroup Path in hierarchy
	Path string

	Resouce *ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path:    path,
	}
}

//add pid process into every cgroup(memory cpu)
func (c *CgroupManager) Apply(pid int) error{
	for _, SubsystemsIns := range(SubsystemsIns) {
		SubsystemsIns.Apply(c.Path, pid)
	}
	return nil
}

//set res for every cgroup in subsystem
func (c *CgroupManager) Set(res *ResourceConfig) error {
	for _,SubsystemsIns := range SubsystemsIns {
		SubsystemsIns.Set(c.Path, res)
	}
	return nil
}

func (c *CgroupManager) Destroy() error {
	for _,SubsystemsIns := range(SubsystemsIns) {
		if err := SubsystemsIns.Remove(c.Path); err != nil {
			log.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}