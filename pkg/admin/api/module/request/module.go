package request

import (
	"main/common/types/http"
	"main/pkg/common/dbm/dao/moduleproperty"
	"time"
)

type ModuleSearch struct {
	ID          uint   `json:"id" form:"id"`
	ParentId    uint   `json:"parentId" form:"parentId"`
	Name        string `json:"name" form:"name"`
	Enable      int    `json:"enable" form:"enable"`
	Description string `json:"description" form:"description"`
	Status      int    `json:"status" form:"status"`
	Plugin      string `json:"plugin" form:"plugin"`
	Type        string `json:"type" form:"type"`

	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	http.PageInfo
}

type ModuleRouterSearch struct {
	Name string `json:"name" form:"name"`
	http.PageInfo
}

type ModuleUpdate struct {
	Id          uint                             `json:"id" form:"id"`
	ParentId    uint                             `json:"parentId" form:"parentId"`
	Name        string                           `json:"name" form:"name"`
	Enable      int                              `json:"enable" form:"enable"`
	Description string                           `json:"description" form:"description"`
	Status      int                              `json:"status" form:"status"`
	Plugin      string                           `json:"plugin" form:"plugin"`
	Type        string                           `json:"type" form:"type"`
	Properties  *[]moduleproperty.ModuleProperty `json:"properties" form:"properties"`
}

type ModuleCreate struct {
	Id          uint                             `json:"id" form:"id"`
	ParentId    uint                             `json:"parentId" form:"parentId"`
	Name        string                           `json:"name" form:"name"`
	Enable      bool                             `json:"enable" form:"enable"`
	Description string                           `json:"description" form:"description"`
	Status      int                              `json:"status" form:"status"`
	Plugin      string                           `json:"plugin" form:"plugin"`
	Type        string                           `json:"type" form:"type"`
	Properties  *[]moduleproperty.ModuleProperty `json:"properties" form:"properties"`
}
