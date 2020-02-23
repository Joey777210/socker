package subsystem

//resource limit config
type ResourceConfig struct{
	MemoryLimit string
	CpuShare string		//cpu time piece weight
	CpuSet string		//cpu core num
}


type Subsystem interface {
	Name() string

	Set (path string, res *ResourceConfig) error

	Apply(path string, pid int) error

	Remove(path string) error
}

var (
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)