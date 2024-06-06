package module

import (
	"fmt"
	"k8s.io/klog/v2"
	"main/common/types/http"
	core "main/pkg/actor/core"
	"main/pkg/admin/api/module/request"
	"main/pkg/admin/api/module/response"
	"main/pkg/common/dbm"
	"main/pkg/common/dbm/dao/common"
	"main/pkg/common/dbm/dao/module"
	"main/pkg/common/dbm/dao/moduleproperty"
	"main/pkg/common/dbm/dbtools"
	"main/pkg/common/modules"
	"main/pkg/north"
)

// CreateModule 创建module表记录
func CreateModule(moduleCreate *request.ModuleCreate) (err error) {
	// check
	exist := module.Exist("name", moduleCreate.Name)
	if exist {
		return fmt.Errorf("%s 已经存在", moduleCreate.Name)
	}

	m := module.Module{
		BasicModel:  common.BasicModel{Id: moduleCreate.Id},
		Name:        moduleCreate.Name,
		Enable:      true,
		Description: "DMP Gateway",
		Plugin:      modules.DmpModuleName,
		Type:        modules.NorthGroup,
		Status:      1,
		ParentId:    moduleCreate.ParentId,
	}
	err = module.AddModuleTrans(&m, moduleCreate.Properties)
	if err == nil {
		north.RegisterModuleFromDB(&m, moduleCreate.Properties)
		err = StartModule(m.Name)
		if err != nil {
			klog.Error("start module error:", err)
		}
	}
	return err
}

// DeleteModule 删除module表记录
func DeleteModule(id uint) (err error) {
	err = module.DeleteByID(dbm.DBAccess, id)
	return err
}

// DeleteModuleByIds 批量删除module表记录
func DeleteModuleByIds(idsReq http.IdsReq) (err error) {
	err = module.DeleteTrans(idsReq.Ids)
	return err
}

// UpdateModule 更新module表记录
func UpdateModule(m request.ModuleUpdate) (err error) {
	params := dbtools.ApplyUpdateParams(&m)
	err = module.UpdateFieldsTrans(m.Id, m.Properties, params)
	return
}

// GetModule 根据id获取module表记录
func GetModule(id uint) (m module.Module, err error) {
	err = module.GetByID(id, &m)
	return
}

// GetModuleAndProperty 根据id获取module and properties
func GetModuleAndProperty(id uint) (*response.ModuleResponse, error) {
	m := module.Module{}
	err := module.GetByID(id, &m)
	if err != nil {
		klog.Error("GetModule error:", err)
		return nil, err
	}
	moduleProperty, err := moduleproperty.QueryModuleProperty(m.Id)
	if err != nil {
		klog.Error("QueryModuleProperty error:", err)
	}
	moduleResponse := response.ModuleResponse{}
	moduleResponse.Module = &m
	moduleResponse.Properties = moduleProperty
	return &moduleResponse, nil
}

// GetModuleByName 根据name获取module表记录
func GetModuleByName(name string) (*response.ModuleResponse, error) {
	m := module.Module{}
	err := module.GetByName(name, &m)
	if err != nil {
		klog.Error("GetModule error:", err)
		return nil, err
	}
	moduleProperty, err := moduleproperty.QueryModuleProperty(m.Id)
	if err != nil {
		klog.Error("QueryModuleProperty error:", err)
	}
	moduleResponse := response.ModuleResponse{}
	moduleResponse.Module = &m
	moduleResponse.Properties = moduleProperty
	return &moduleResponse, nil
}

// GetModuleInfoList 分页获取module表记录
func GetModuleInfoList(info *request.ModuleSearch) (list *[]module.Module, total int64, err error) {
	return module.GetInfoList(info)
}

// GetModuleListWithProperty 分页获取module表记录
func GetModuleListWithProperty(info *request.ModuleSearch) (*[]response.DmpGatewaySubDeviceModuleResponse, int64, error) {
	moduleInfoList, total, err := module.GetInfoList(info)
	var subDeviceModuleResponses []response.DmpGatewaySubDeviceModuleResponse
	if err != nil || moduleInfoList == nil || len(*moduleInfoList) == 0 {
		return &subDeviceModuleResponses, 0, err
	}
	for _, m := range *moduleInfoList {
		moduleResponse := response.DmpGatewaySubDeviceModuleResponse{Module: m}
		moduleProperty, err := moduleproperty.QueryModuleProperty(m.Id)
		if err != nil {
			klog.Error("QueryModuleProperty error:", err)
			continue
		}
		for _, p := range *moduleProperty {
			if p.Name == "productKey" {
				moduleResponse.ProductKey = p.Value
				continue
			} else if p.Name == "deviceKey" {
				moduleResponse.DeviceKey = p.Value
				continue
			} else if p.Name == "deviceSecret" {
				moduleResponse.DeviceSecret = p.Value
				continue
			} else if p.Name == "productName" {
				moduleResponse.ProductName = p.Value
				continue
			} else if p.Name == "deviceName" {
				moduleResponse.DeviceName = p.Value
				continue
			}
		}
		subDeviceModuleResponses = append(subDeviceModuleResponses, moduleResponse)
	}
	return &subDeviceModuleResponses, total, nil
}

// GetModuleInfoListRouter 分页获取module表记录
func GetModuleInfoListRouter(info *request.ModuleSearch) (list *[]response.ModuleRouterResponse, total int64, err error) {
	modules, total, err := module.GetInfoListInRouter(info)
	var subDeviceModuleResponses []response.ModuleRouterResponse
	if err != nil || modules == nil || len(*modules) == 0 {
		return nil, 0, err
	}
	for _, mod := range *modules {
		var resp response.ModuleRouterResponse
		moduleProperty, pErr := moduleproperty.QueryModuleProperty(mod.Id)
		if pErr != nil {
			klog.Error("QueryModuleProperty error:", err)
			continue
		}
		productKey := "-"
		deviceKey := "-"
		for _, p := range *moduleProperty {
			if p.Name == "productKey" {
				productKey = p.Value
				continue
			} else if p.Name == "deviceKey" {
				deviceKey = p.Value
				continue
			}
		}
		resp.DeviceKey = deviceKey
		resp.Topic = "/GEWU/" + productKey + "/" + deviceKey
		mod.Name, err = GetModuleName(mod)
		resp.Id = mod.Id
		resp.RouterPlace = mod.Name
		subDeviceModuleResponses = append(subDeviceModuleResponses, resp)
	}
	return &subDeviceModuleResponses, total, nil
}

func GetModuleName(mod module.Module) (string, error) {
	productName, productNameErr := moduleproperty.GetByModuleId(mod.Id, "productName")
	deviceName, deviceNameErr := moduleproperty.GetByModuleId(mod.Id, "deviceName")
	if productNameErr != nil || deviceNameErr != nil {
		return "GEWU", productNameErr
	}
	return "GEWU/" + productName.Value + "/" + deviceName.Value, nil

}

func AddModuleTrans(add *module.Module, addAttrs *[]moduleproperty.ModuleProperty) error {
	err := module.AddModuleTrans(add, addAttrs)
	if err != nil {
		return err
	}
	return nil
}

func StartModule(name string) error {
	moduleMap := core.GetModules()
	moduleInfo := moduleMap[name]
	moduleInfo.GetModule().SetEnable(true)
	core.StartModule(name, moduleMap[name])
	return nil
}

func StopModule(name string) error {
	moduleMap := core.GetModules()
	moduleInfo := moduleMap[name]
	moduleInfo.GetModule().Stop()
	return nil
}
