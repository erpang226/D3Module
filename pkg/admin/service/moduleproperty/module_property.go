package moduleproperty

import (
	"fmt"
	"main/common/tools"
	"main/common/types/http"
	"main/pkg/admin/api/moduleproperty/request"
	"main/pkg/common/dbm"
	"main/pkg/common/dbm/dao/moduleproperty"
)

// CreatModuleProperty 创建module表记录
func CreatModuleProperty(m *moduleproperty.ModuleProperty) (err error) {
	// check
	exist := moduleproperty.Exist(m.Name, m.ModuleId)
	if exist {
		return fmt.Errorf("%s 已经存在", m.Name)
	}
	now := tools.GetNow()
	m.CreatedAt = now
	m.UpdatedAt = now
	err = moduleproperty.Save(dbm.DBAccess, m)
	return err
}

// DeleteModuleProperty 删除module表记录
func DeleteModuleProperty(id uint) (err error) {
	err = moduleproperty.DeleteById(id)
	return err
}

// DeleteModulePropertyByIds 批量删除module表记录
func DeleteModulePropertyByIds(idsReq http.IdsReq) (err error) {
	err = moduleproperty.DeleteTrans(idsReq.Ids)
	return err
}

// UpdateModuleProperty 更新module表记录
func UpdateModuleProperty(m request.ModulePropertyUpdate) (err error) {
	err = moduleproperty.UpdateFields(dbm.DBAccess, m)
	return
}

// GetModuleProperty 根据id获取module表记录
func GetModuleProperty(id uint) (m moduleproperty.ModuleProperty, err error) {
	err = moduleproperty.GetByID(id, &m)
	return
}

// GetModulePropertyInfoList 分页获取module表记录
func GetModulePropertyInfoList(info *request.ModulePropertySearch) (list *[]moduleproperty.ModuleProperty, total int64, err error) {
	return moduleproperty.GetInfoList(info)
}
