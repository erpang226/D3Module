package plugin

import (
	"main/pkg/admin/api/plugin/request"
	"main/pkg/common/dbm"
	"main/pkg/common/dbm/dao/common"
	"main/pkg/common/dbm/dbtools"

	"k8s.io/klog/v2"
)

type Plugin struct {
	common.BasicModel
	Name        string `orm:"column(name); not null" json:"name" form:"name"`
	Description string `orm:"column(description); null" json:"description" form:"description"`
	NodeType    int    `orm:"column(node_type); null" json:"nodeType" form:"nodeType"`
	Type        int    `orm:"column(type); null" json:"type" form:"type"`
}

func Save(doc *Plugin) error {
	num, err := dbm.DBAccess.Insert(doc)
	klog.V(4).Infof("Insert affected Num: %d", num)
	return err
}

func DeleteByID(id uint) error {
	filtered := dbm.DBAccess.QueryTable(common.PluginTableName).Filter("id", id)
	num, err := filtered.Delete()
	if err != nil {
		klog.Errorf("Something wrong when deleting data: %v", err)
		return err
	}
	klog.V(4).Infof("Delete affected Num: %d", num)
	return nil
}

func UpdateField(deviceID string, col string, value interface{}) error {
	_, err := dbm.DBAccess.QueryTable(common.PluginTableName).Filter("id", deviceID).Update(map[string]interface{}{col: value})
	return err
}

func UpdateFields(deviceID uint, cols map[string]interface{}) error {
	filtered := dbm.DBAccess.QueryTable(common.PluginTableName).Filter("id", deviceID)
	num, err := filtered.Update(cols)
	klog.V(4).Infof("Update affected Num: %d", num)
	return err
}

func DeleteTrans(deletes []uint) error {
	obm := dbm.DefaultOrmFunc()
	err := obm.Begin()
	if err != nil {
		klog.Errorf("failed to begin transaction: %v", err)
		return err
	}

	defer func() {
		if err != nil {
			dbm.RollbackTransaction(obm)
		} else {
			err = obm.Commit()
			if err != nil {
				klog.Errorf("failed to commit transaction: %v", err)
			}
		}
	}()

	for _, del := range deletes {
		err = DeleteByID(del)
		if err != nil {
			return err
		}
	}

	return err
}

func Query(key string, condition string) (*[]Plugin, error) {
	devices := new([]Plugin)
	_, err := dbm.DBAccess.QueryTable(common.PluginTableName).Filter(key, condition).All(devices)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func QueryDeviceAll() (*[]Plugin, error) {
	devices := new([]Plugin)
	_, err := dbm.DBAccess.QueryTable(common.PluginTableName).All(devices)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

type PluginUpdate struct {
	Id   uint `json:"id" form:"id"`
	Cols map[string]interface{}
}

func UpdateDeviceMulti(updates []PluginUpdate) error {
	var err error
	for _, update := range updates {
		err = UpdateFields(update.Id, update.Cols)
		if err != nil {
			return err
		}
	}
	return nil
}

func Exist(key string, condition any) bool {
	exist := dbm.DBAccess.QueryTable(common.PluginTableName).Filter(key, condition).Exist()
	return exist
}

func GetByID(id uint, row *Plugin) error {
	err := dbm.DBAccess.QueryTable(common.PluginTableName).Filter("id", id).One(row)
	if err != nil {
		klog.Errorf("Something wrong when get data: %v", err)
		return err
	}
	return nil
}

func GetInfoList(info *request.PluginSearch) (*[]Plugin, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	deviceModelList := new([]Plugin)
	cond := dbtools.ApplyCondition(info)
	total, err := dbm.DBAccess.QueryTable(common.PluginTableName).SetCond(cond).Count()
	if err != nil {
		klog.Error("分页获取device model表记录", err)
		return nil, total, err
	}
	_, err1 := dbm.DBAccess.QueryTable(common.PluginTableName).SetCond(cond).Limit(limit).Offset(offset).All(deviceModelList)
	if err1 != nil {
		klog.Error("分页获取device model表记录", err)
		return nil, total, err
	}
	return deviceModelList, total, err
}
