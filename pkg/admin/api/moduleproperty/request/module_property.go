package request

import (
	"main/common/types/http"
	"time"
)

type ModulePropertySearch struct {
	ID          uint   `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	ModuleId    uint   `json:"moduleId" form:"moduleId"`
	Description string `json:"description" form:"description"`

	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	http.PageInfo
}

type ModulePropertyUpdate struct {
	Id       uint                   `json:"id" form:"id"`
	ModuleId uint                   `json:"moduleId" form:"moduleId"`
	Name     string                 `json:"name" form:"name"`
	Cols     map[string]interface{} `json:"cols" form:"cols"`
}
