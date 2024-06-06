package core

import (
	"main/pkg/common/dbm/dao/module"
	"main/pkg/common/dbm/dao/moduleproperty"
	"main/pkg/common/dbm/dao/plugin"

	"github.com/beego/beego/orm"
	"k8s.io/klog/v2"
)

// RegisterCoreDBTable create table
func RegisterCoreDBTable() {
	klog.Infof("Begin to register core db model")

	orm.RegisterModel(new(module.Module), new(moduleproperty.ModuleProperty))
	orm.RegisterModel(new(plugin.Plugin))
}
