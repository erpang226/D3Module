package moduleproperty

import (
	"fmt"
	"github.com/beego/beego/orm"
	"k8s.io/klog/v2"
	"main/common/tools"
	"main/pkg/admin/api/moduleproperty/request"
	"main/pkg/common/dbm"
	"main/pkg/common/dbm/dao/common"
	"main/pkg/common/dbm/dbtools"
)

type ModuleProperty struct {
	common.BasicModel
	Name        string `orm:"column(name);not null"  json:"name" form:"name"`
	ModuleId    uint   `orm:"column(module_id);not null" json:"moduleId" form:"moduleId"`
	DataType    string `orm:"column(data_type);not null" json:"dataType" form:"dataType"`
	Description string `orm:"column(description);null" json:"description" form:"description"`
	Value       string `orm:"column(value);null" json:"value" form:"value"`
}

// Save module property
func Save(obm orm.Ormer, doc *ModuleProperty) error {
	now := tools.GetNow()
	doc.CreatedAt = now
	doc.UpdatedAt = now
	_, err := obm.Insert(doc)
	return err
}

// Exist 是否已存在
func Exist(name string, moduleId uint) bool {
	exist := dbm.DBAccess.QueryTable(common.ModulePropertyTableName).Filter("name", name).Filter("module_id", moduleId).Exist()
	return exist
}

func GetByModuleId(moduleId uint, name string) (*ModuleProperty, error) {
	var doc ModuleProperty
	err := dbm.DBAccess.QueryTable(common.ModulePropertyTableName).Filter("name", name).Filter("module_id", moduleId).One(&doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil

}

// DeleteByModuleID delete module property
func DeleteByModuleID(obm orm.Ormer, moduleId uint) error {
	num, err := obm.QueryTable(common.ModulePropertyTableName).Filter("module_id", moduleId).Delete()
	if err != nil {
		klog.Error("Something wrong when deleting module property", err)
		return err
	}
	klog.V(4).Info("Delete module property affected Num: ", num)
	return nil
}

// DeleteById delete module property
func DeleteById(id uint) error {
	num, err := dbm.DBAccess.QueryTable(common.ModulePropertyTableName).Filter("id", id).Delete()
	if err != nil {
		klog.Errorf("Something wrong when deleting module property: %v", err)
		return err
	}
	klog.V(4).Infof("Delete module property affected Num: %d", num)
	return nil
}

// DeleteTrans the transaction of delete Module
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
		err = DeleteById(del)
		if err != nil {
			return err
		}
	}

	return err
}

// DeleteModuleProperty delete module property
func DeleteModuleProperty(obm orm.Ormer, moduleId uint, propertyName string) error {
	num, err := obm.QueryTable(common.ModulePropertyTableName).Filter("module_id", moduleId).Filter("name", propertyName).Delete()
	if err != nil {
		klog.Errorf("Something wrong when deleting module property: %v", err)
		return err
	}
	klog.V(4).Infof("Delete module property affected Num: %d", num)
	return nil
}

// UpdateField update special field
func UpdateField(moduleId uint, propertyName string, col string, value interface{}) error {
	num, err := dbm.DBAccess.QueryTable(common.ModulePropertyTableName).Filter("module_id", moduleId).Filter("name", propertyName).Update(map[string]interface{}{col: value})
	klog.V(4).Infof("Update affected Num: %d", num)
	return err
}

// UpdateFields update special fields
func UpdateFields(obm orm.Ormer, model request.ModulePropertyUpdate) error {
	params := dbtools.ApplyUpdateParams(&model)
	num, err := obm.QueryTable(common.ModulePropertyTableName).Filter("id", model.Id).Update(params)
	klog.V(4).Infof("Update affected Num: %d", num)
	return err
}

// GetInfoList 分页获取module表记录
func GetInfoList(info *request.ModulePropertySearch) (*[]ModuleProperty, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	moduleList := new([]ModuleProperty)
	cond := dbtools.ApplyCondition(info)
	total, err := dbm.DBAccess.QueryTable(common.ModulePropertyTableName).SetCond(cond).Count()
	if err != nil {
		klog.Error("分页获取module表记录", err)
		return nil, total, err
	}
	_, err1 := dbm.DBAccess.QueryTable(common.ModulePropertyTableName).SetCond(cond).Limit(limit).Offset(offset).All(moduleList)
	if err1 != nil {
		klog.Error("分页获取module表记录", err)
		return nil, total, err
	}
	return moduleList, total, err
}

// GetByID get Module by id
func GetByID(id uint, row *ModuleProperty) error {
	err := dbm.DBAccess.QueryTable(common.ModulePropertyTableName).Filter("id", id).One(row)
	if err != nil {
		klog.Errorf("Something wrong when get data: %v", err)
		return err
	}
	return nil
}

// QueryModuleProperty query ModuleProperty
func QueryModuleProperty(moduleId uint) (*[]ModuleProperty, error) {
	attrs := new([]ModuleProperty)
	_, err := dbm.DBAccess.QueryTable(common.ModulePropertyTableName).Filter("module_id", moduleId).All(attrs)
	if err != nil {
		return nil, err
	}
	return attrs, nil
}

// QueryModulePropertyInModuleIds query ModuleProperty
func QueryModulePropertyInModuleIds(moduleIds []uint) (*[]ModuleProperty, error) {
	attrs := new([]ModuleProperty)
	cond := orm.NewCondition()
	cond = cond.And(fmt.Sprintf("%s__%s", "module_id", "in"), moduleIds)
	_, err := dbm.DBAccess.QueryTable(common.ModulePropertyTableName).SetCond(cond).All(attrs)
	if err != nil {
		return nil, err
	}
	return attrs, nil
}

// ModulePropertyDelete the struct for deleting device
type ModulePropertyDelete struct {
	ModuleId uint
	Name     string
}

// ModulePropertyTrans transaction of module property
func ModulePropertyTrans(adds []ModuleProperty, deletes []ModulePropertyDelete, updates []request.ModulePropertyUpdate) error {
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

	for _, add := range adds {
		err = Save(obm, &add)
		if err != nil {
			return err
		}
	}

	for _, del := range deletes {
		err = DeleteModuleProperty(obm, del.ModuleId, del.Name)
		if err != nil {
			return err
		}
	}

	for _, update := range updates {
		err = UpdateFields(obm, update)
		if err != nil {
			return err
		}
	}

	return err
}
