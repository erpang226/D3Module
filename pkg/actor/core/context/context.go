package context

import (
	"time"

	"main/pkg/actor/common"
	"main/pkg/actor/core/model"
)

// ModuleContext is interface for context module management
type ModuleContext interface {
	AddModule(info *common.ModuleInfo)
	AddModuleGroup(module, group string)
	Cleanup(module string)
}

// MessageContext is interface for message syncing
type MessageContext interface {
	// async mode
	Send(module string, message model.Message)
	Receive(module string) (model.Message, error)
	// sync mode
	SendSync(module string, message model.Message, timeout time.Duration) (model.Message, error)
	SendResp(message model.Message)
	// group broadcast
	SendToGroup(group string, message model.Message)
	SendToGroupSync(group string, message model.Message, timeout time.Duration) error
}
