package module

import (
	"github.com/beego/beego/orm"
	"k8s.io/klog/v2"
	"main/common/tools"
	"main/pkg/admin/api/module/request"
	"main/pkg/common/dbm"
	"main/pkg/common/dbm/dao/common"
	"main/pkg/common/dbm/dao/moduleproperty"
	"main/pkg/common/dbm/dbtools"
	"main/pkg/common/modules"
)

type Module struct {
	common.BasicModel
	Name        string `orm:"column(name); not null" json:"name" form:"name"`
	Enable      bool   `orm:"column(enable); null" json:"enable" form:"enable"`
	Description string `orm:"column(description); null" json:"description" form:"description"`
	Status      int    `orm:"column(status); null" json:"status" form:"status"`
	Plugin      string `orm:"column(plugin); not null" json:"plugin" form:"plugin"`
	Type        string `orm:"column(type); not null" json:"type" form:"type"`
	ParentId    uint   `orm:"column(parent_id); not null" json:"parentId" form:"parentId"`
}

// Save Module
func Save(obm orm.Ormer, doc *Module) error {
	num, err := obm.Insert(doc)
	klog.V(4).Infof("Insert affected Num: %d", num)
	return err
}

// DeleteByID delete Module by id
func DeleteByID(obm orm.Ormer, id uint) error {
	num, err := obm.QueryTable(common.ModuleTableName).Filter("id", id).Delete()
	if err != nil {
		klog.Errorf("Something wrong when deleting data: %v", err)
		return err
	}
	klog.V(4).Infof("Delete affected Num: %d", num)
	return nil
}

// UpdateField update special field
func UpdateField(ModuleID string, col string, value interface{}) error {
	num, err := dbm.DBAccess.QueryTable(common.ModuleTableName).Filter("id", ModuleID).Update(map[string]interface{}{col: value})
	klog.V(4).Infof("Update affected Num: %d", num)
	return err
}

// UpdateFields update special fields
func UpdateFields(model request.ModuleUpdate) error {
	params := dbtools.ApplyUpdateParams(&model)
	num, err := dbm.DBAccess.QueryTable(common.ModuleTableName).Filter("id", model.Id).Update(params)
	klog.V(4).Infof("Update affected Num: %d", num)
	return err
}

// GetInfoList 分页获取module表记录
func GetInfoList(info *request.ModuleSearch) (*[]Module, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	moduleList := new([]Module)
	cond := dbtools.ApplyCondition(info)
	// todo delete
	if info.ParentId == 0 {
		cond = cond.And("plugin__ne", modules.DmpModuleName)
	}
	total, err := dbm.DBAccess.QueryTable(common.ModuleTableName).SetCond(cond).Count()
	if err != nil {
		klog.Error("分页获取module表记录", err)
		return nil, total, err
	}
	_, err1 := dbm.DBAccess.QueryTable(common.ModuleTableName).SetCond(cond).Limit(limit).Offset(offset).All(moduleList)
	if err1 != nil {
		klog.Error("分页获取module表记录", err)
		return nil, total, err
	}
	return moduleList, total, err
}

// GetInfoListInRouter 分页获取module表记录
func GetInfoListInRouter(info *request.ModuleSearch) (*[]Module, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	moduleList := new([]Module)
	cond := dbtools.ApplyCondition(info)
	// todo delete
	if info.ParentId == 0 {
		cond = cond.And("name__ne", "GEWUDMP")
	}
	total, err := dbm.DBAccess.QueryTable(common.ModuleTableName).SetCond(cond).Count()
	if err != nil {
		klog.Error("分页获取module表记录", err)
		return nil, total, err
	}
	_, err1 := dbm.DBAccess.QueryTable(common.ModuleTableName).SetCond(cond).Limit(limit).Offset(offset).All(moduleList)
	if err1 != nil {
		klog.Error("分页获取module表记录", err)
		return nil, total, err
	}
	return moduleList, total, err
}

// GetByID get Module by id
func GetByID(id uint, row *Module) error {
	err := dbm.DBAccess.QueryTable(common.ModuleTableName).Filter("id", id).One(row)
	if err != nil {
		klog.Errorf("Something wrong when get data: %v", err)
		return err
	}
	return nil
}

// GetByName get Module by name
func GetByName(name string, row *Module) error {
	err := dbm.DBAccess.QueryTable(common.ModuleTableName).Filter("name", name).One(row)
	if err != nil {
		klog.Errorf("Something wrong when get data: %v", err)
		return err
	}
	return nil
}

// Exist 是否已存在
func Exist(key string, condition string) bool {
	exist := dbm.DBAccess.QueryTable(common.ModuleTableName).Filter(key, condition).Exist()
	return exist
}

// Query query Module
func Query(key string, condition string) (*[]Module, error) {
	Modules := new([]Module)
	_, err := dbm.DBAccess.QueryTable(common.ModuleTableName).Filter(key, condition).All(Modules)
	if err != nil {
		return nil, err
	}
	return Modules, nil
}

// QueryAll query all
func QueryAll() (*[]Module, error) {
	Modules := new([]Module)
	_, err := dbm.DBAccess.QueryTable(common.ModuleTableName).All(Modules)
	if err != nil {
		return nil, err
	}
	return Modules, nil
}

// AddModuleTrans the transaction of add Module
func AddModuleTrans(add *Module, addAttrs *[]moduleproperty.ModuleProperty) error {
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
	now := tools.GetNow()
	add.CreatedAt = now
	add.UpdatedAt = now
	err = Save(obm, add)
	if err != nil {
		klog.Errorf("save Module failed: %v", err)
		return err
	}
	for _, attr := range *addAttrs {
		attr.ModuleId = add.Id
		err = moduleproperty.Save(obm, &attr)
		if err != nil {
			return err
		}
	}

	return err
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
		err = DeleteByID(obm, del)
		if err != nil {
			return err
		}
		err = moduleproperty.DeleteByModuleID(obm, del)
		if err != nil {
			return err
		}
		//级联删除通道？

	}

	return err
}

func UpdateFieldsTrans(moduleId uint, properties *[]moduleproperty.ModuleProperty, cols map[string]interface{}) error {
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
	_, moduleErr := dbm.DBAccess.QueryTable(common.ModuleTableName).Filter("id", moduleId).Update(cols)
	if moduleErr != nil {
		err = moduleErr
		return err
	}
	for _, property := range *properties {
		property.ModuleId = moduleId
		params := dbtools.ApplyUpdateParams(&property)
		_, propertyErr := obm.QueryTable(common.ModulePropertyTableName).Filter("module_id", moduleId).Filter("name", property.Name).Update(params)
		if propertyErr != nil {
			klog.Errorf("failed to update property: %v", propertyErr)
			err = propertyErr
			return err
		}
	}
	return err

}
