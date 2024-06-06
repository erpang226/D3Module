package adapter

import (
	"main/pkg/actor/core/model"
)

// Adapter is a south client interface
type Adapter interface {
	Init() error
	UnInit()
	SendTOSouth(message model.Message, params ...string) error
	ReceiveFromSouth(params ...string) (model.Message, error)
}
