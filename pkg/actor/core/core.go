package actor

import (
	"main/pkg/actor/common"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/klog/v2"

	actorContext "main/pkg/actor/core/context"
)

func StartModules() {
	actorContext.InitContext([]string{common.MsgCtxTypeChannel})

	modules := GetModules()

	for name, module := range modules {
		StartModule(name, module)
	}
}

func StartModule(name string, module *ModuleInfo) {
	var m common.ModuleInfo
	switch module.contextType {
	case common.MsgCtxTypeChannel:
		m = common.ModuleInfo{
			ModuleName: name,
			ModuleType: module.contextType,
		}
	case common.MsgCtxTypeUS:
		m = common.ModuleInfo{
			ModuleName: name,
			ModuleType: module.contextType,
			ModuleSocket: common.ModuleSocket{
				IsRemote: module.remote,
			},
		}
	default:
		klog.Exitf("unsupported context type: %s", module.contextType)
	}

	actorContext.AddModule(&m)
	actorContext.AddModuleGroup(name, module.module.Group())

	go moduleKeeper(name, module, m)
	klog.Infof("starting module %s", name)
}

func GracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT)
	s := <-c
	klog.Infof("Get os signal %v", s.String())

	actorContext.Cancel()
	modules := GetModules()
	for name := range modules {
		klog.Infof("Cleanup module %v", name)
		actorContext.Cleanup(name)
	}
}

func Run() {
	StartModules()
	GracefulShutdown()
}

func moduleKeeper(name string, moduleInfo *ModuleInfo, m common.ModuleInfo) {
	for {
		moduleInfo.module.Start()
		if !moduleInfo.remote {
			return
		}
		actorContext.AddModule(&m)
		actorContext.AddModuleGroup(name, moduleInfo.module.Group())
	}
}
