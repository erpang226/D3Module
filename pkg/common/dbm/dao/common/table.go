package common

const (
	//DeviceModelTableName device table
	DeviceModelTableName = "device_model"
	//DeviceModelPropertyTableName device property table
	DeviceModelPropertyTableName = "device_model_property"
	//ModuleTableName module table
	ModuleTableName = "module"
	//ModulePropertyTableName module property table
	ModulePropertyTableName    = "module_property"
	PointGroupTableName        = "point_group"
	PointTableName             = "point"
	RouterTableName            = "router"
	RouterConfigTableName      = "router_config"
	EquipmentTableName         = "equipment"
	EquipmentPropertyTableName = "equipment_property"
	EquipmentChannelTableName  = "equipment_channel"
	EquipmentGroupTableName    = "equipment_group"
	PluginTableName            = "plugin"
)

type BasicModel struct {
	Id        uint   `orm:"primarykey" json:"id" form:"id"`                       // 主键ID
	CreatedAt string `orm:"column(created_at)" json:"createdAt" form:"createdAt"` // 创建时间
	UpdatedAt string `orm:"column(updated_at)" json:"updatedAt" form:"updatedAt"` // 更新时间
}
