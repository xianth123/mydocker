package subsystems

type ResourceCofing struct {
	MemoryLimit	string
	CpuShare	string
	CpuSet		string
}

type Subsystem interface {
	Name() string 				// 返回资源的名字
	Set(path string, res *ResourceCofing) error // 使用 ResourceConfig 来设置对应的 subsystem
	Apply(path string, pid int) error		// 将 pid 加入到 subsystem
	Remove(path string) error				// 从 subsystem 中移除 pid
}

// 每个 subsystem 都代表一种资源的描述
// 进程对资源的描述是 三种不同资源描述指针组成的列表
// 3 个不同子系统叠加描述了一个 docker 的资源限制
var (
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},

	}
)
