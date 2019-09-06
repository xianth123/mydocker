package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CpuSubSystem struct {

}

func (c *CpuSubSystem) Set(cgroupPath string, res *ResourceCofing) error{
	if subsysCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, true); err == nil {
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(res.CpuShare), 0644); err != nil {
			return fmt.Errorf("set cgroup memory fail %v", err)
		}
		return nil
	} else {
		return err
	}
}

func (c *CpuSubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, true); err == nil {
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}

func (c *CpuSubSystem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(c.Name(), cgroupPath, true); err == nil {
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	} else {
		return err
	}
}

func (c *CpuSubSystem) Name() string  {
	return "cpu"
}
