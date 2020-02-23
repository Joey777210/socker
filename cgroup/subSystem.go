package cgroup

//resource limit config
type ResourceConfig struct{
	MemoryLimit string
	CpuShare string		//cpu time piece weight
	CpuSet string		//cpu core num
}


type Subsystem interface {
	//return the name of current subsystem
	Name() string
	//set res to cgroup
	Set (path string, res *ResourceConfig) error
	//apply cgroup to pid process || add pid to cgroup
	Apply(path string, pid int) error
	//remove a cgroup in parmPath
	Remove(path string) error
}

var (
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)