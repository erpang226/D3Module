package south

import (
	"main/pkg/common/dbm/dao/module"
	"main/pkg/common/dbm/dao/moduleproperty"

	"k8s.io/klog/v2"
)

func RegisterModuleFromDB(m *module.Module, property *[]moduleproperty.ModuleProperty) {
	switch m.Plugin {
	default:
		klog.Errorf("module %s error plugin. %s", m.Name, m.Plugin)
	}
}
