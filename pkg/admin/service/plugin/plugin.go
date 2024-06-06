package plugin

import (
	"fmt"
	"main/common/tools"
	"main/common/types/http"
	"main/pkg/admin/api/plugin/request"
	"main/pkg/common/dbm/dao/plugin"
	"main/pkg/common/dbm/dbtools"
)

// CreatePlugin 创建deviceModel表记录
func CreatePlugin(m *plugin.Plugin) (err error) {
	// check
	exist := plugin.Exist("name", m.Name)
	if exist {
		return fmt.Errorf("%s 已经存在", m.Name)
	}
	now := tools.GetNow()
	m.CreatedAt = now
	m.UpdatedAt = now
	err = plugin.Save(m)
	return err
}

// DeletePluginById 删除deviceModel表记录
func DeletePluginById(id uint) (err error) {
	err = plugin.DeleteByID(id)
	return err
}

// DeletePluginByIds 批量删除deviceModel表记录
func DeletePluginByIds(ids http.IdsReq) (err error) {
	err = plugin.DeleteTrans(ids.Ids)
	return err
}

// UpdatePlugin 更新deviceModel表记录
func UpdatePlugin(m request.PluginUpdate) (err error) {
	params := dbtools.ApplyUpdateParams(&m)
	err = plugin.UpdateFields(m.Id, params)
	return err
}

// GetPlugin 根据id获取deviceModel表记录
func GetPlugin(id uint) (m plugin.Plugin, err error) {
	err = plugin.GetByID(id, &m)
	return m, err
}

// GetPluginInfoList 分页获取deviceModel表记录
func GetPluginInfoList(info *request.PluginSearch) (list *[]plugin.Plugin, total int64, err error) {
	return plugin.GetInfoList(info)
}
