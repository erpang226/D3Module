package response

import (
	"main/pkg/common/dbm/dao/module"
	"main/pkg/common/dbm/dao/moduleproperty"
)

type ModuleResponse struct {
	*module.Module
	Properties *[]moduleproperty.ModuleProperty `json:"properties"`
}

type ModuleRouterResponse struct {
	Id          uint   `json:"id"`
	DeviceKey   string `json:"deviceKey"`
	RouterPlace string `json:"routerPlace"`
	Topic       string `json:"topic"`
}

type DmpGatewaySubDeviceModuleResponse struct {
	module.Module
	ProductKey   string `json:"productKey"`
	DeviceKey    string `json:"deviceKey"`
	DeviceSecret string `json:"deviceSecret"`
	ProductName  string `json:"productName"`
	DeviceName   string `json:"deviceName"`
}
