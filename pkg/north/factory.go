package north

import (
	"main/pkg/common/dbm/dao/module"
	"main/pkg/common/dbm/dao/moduleproperty"
	"main/pkg/common/modules"
	mqttModule "main/pkg/north/mqtt"

	"k8s.io/klog/v2"
)

func RegisterModuleFromDB(m *module.Module, property *[]moduleproperty.ModuleProperty) {
	switch m.Plugin {
	case modules.MqttModuleName:
		newMqttModule := mqttModule.NewMqttModuleFromDB(m, property)
		mqttModule.Register(newMqttModule)

		//...

	default:
		klog.Errorf("module %s error plugin. %s", m.Name, m.Plugin)
	}
}
