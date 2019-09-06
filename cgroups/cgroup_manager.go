package cgroups

import (
	"github.com/sirupsen/logrus"
	"mydocker/cgroups/subsystems"
	)

type CgroupManager struct {
	Path 	string
	Resource *subsystems.ResourceCofing
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

// SubsystemsIns 是一个初始化的设置的变量
func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range(subsystems.SubsystemsIns) {
		subSysIns.Apply(c.Path, pid)
	}
	return nil
}

// 设置 cgroup 资源限制
func (c *CgroupManager) Set(res *subsystems.ResourceCofing) error {
	for _, subSysIns := range (subsystems.SubsystemsIns) {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

// 释放 cgroup
func (c *CgroupManager) Destroy() error {
	for _, subSysIns := range(subsystems.SubsystemsIns) {
		if err := subSysIns.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}


