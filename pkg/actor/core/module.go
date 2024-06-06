package actor

import (
	"main/pkg/actor/common"
	"main/pkg/actor/core/socket"

	"k8s.io/klog/v2"
)

type Module interface {
	Name() string
	Group() string
	Start()
	Enable() bool
	SetEnable(enable bool)
	Stop()
}

var (
	modules         map[string]*ModuleInfo
	disabledModules map[string]*ModuleInfo
)

func init() {
	modules = make(map[string]*ModuleInfo)
	disabledModules = make(map[string]*ModuleInfo)
}

type ModuleInfo struct {
	contextType string
	remote      bool
	module      Module
}

func Register(m Module, opts ...string) {
	info := &ModuleInfo{
		module:      m,
		contextType: common.MsgCtxTypeChannel,
		remote:      false,
	}

	if len(opts) > 0 {
		info.contextType = opts[0]
		info.remote = true
	}

	if m.Enable() {
		modules[m.Name()] = info
		klog.Infof("Module %s registered successfully", m.Name())
	} else {
		disabledModules[m.Name()] = info
		klog.Warningf("Module %v is disabled, do not register", m.Name())
	}
}

func GetModules() map[string]*ModuleInfo {
	return modules
}

func (m *ModuleInfo) GetModule() Module {
	return m.module
}

func GetModuleExchange() *socket.ModuleExchange {
	exchange := socket.ModuleExchange{
		Groups: make(map[string][]string),
	}
	for name, moduleInfo := range modules {
		exchange.Modules = append(exchange.Modules, name)
		group := moduleInfo.module.Group()
		exchange.Groups[group] = append(exchange.Groups[group], name)
	}
	return &exchange
}
