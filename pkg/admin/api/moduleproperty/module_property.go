package moduleproperty

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	"main/common/types/http"
	"main/pkg/admin/api/moduleproperty/request"
	service "main/pkg/admin/service/moduleproperty"
	dao "main/pkg/common/dbm/dao/moduleproperty"
)

// CreateModuleProperty 创建moduleProperty表
// @Tags ModuleProperty
// @Summary 创建moduleProperty表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dao.ModuleProperty true "创建moduleProperty表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /module/createModuleProperty [post]
func CreateModuleProperty(c *gin.Context) {
	var moduleProperty dao.ModuleProperty
	err := c.ShouldBindJSON(&moduleProperty)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}

	if err = service.CreatModuleProperty(&moduleProperty); err != nil {
		klog.Error("创建失败!", err)
		http.FailWithMessage(err.Error(), c)
	} else {
		http.OkWithMessage("创建成功", c)
	}
}

// DeleteModuleProperty 删除moduleProperty表
// @Tags ModuleProperty
// @Summary 删除moduleProperty表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dao.ModuleProperty true "删除moduleProperty表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /module/deleteModuleProperty [delete]
func DeleteModuleProperty(c *gin.Context) {
	var moduleProperty dao.ModuleProperty
	err := c.ShouldBindJSON(&moduleProperty)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.DeleteModuleProperty(moduleProperty.Id); err != nil {
		klog.Error("删除失败!", err)
		http.FailWithMessage("删除失败", c)
	} else {
		http.OkWithMessage("删除成功", c)
	}
}

// DeleteModulePropertyByIds 批量删除moduleProperty表
// @Tags ModuleProperty
// @Summary 批量删除moduleProperty表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除moduleProperty表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /module/deleteModulePropertyByIds [delete]
func DeleteModulePropertyByIds(c *gin.Context) {
	var IDS http.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.DeleteModulePropertyByIds(IDS); err != nil {
		klog.Error("批量删除失败!", err)
		http.FailWithMessage("批量删除失败", c)
	} else {
		http.OkWithMessage("批量删除成功", c)
	}
}

// UpdateModuleProperty 更新moduleProperty表
// @Tags ModuleProperty
// @Summary 更新moduleProperty表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.ModulePropertyUpdate true "更新moduleProperty表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /module/updateModuleProperty [put]
func UpdateModuleProperty(c *gin.Context) {
	var moduleProperty request.ModulePropertyUpdate
	err := c.ShouldBindJSON(&moduleProperty)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.UpdateModuleProperty(moduleProperty); err != nil {
		klog.Error("更新失败!", err)
		http.FailWithMessage("更新失败", c)
	} else {
		http.OkWithMessage("更新成功", c)
	}
}

// FindModuleProperty 用id查询moduleProperty表
// @Tags ModuleProperty
// @Summary 用id查询moduleProperty表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query dao.ModuleProperty true "用id查询moduleProperty表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /module/findModuleProperty [get]
func FindModuleProperty(c *gin.Context) {
	var moduleProperty dao.ModuleProperty
	err := c.ShouldBindQuery(&moduleProperty)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if property, err := service.GetModuleProperty(moduleProperty.Id); err != nil {
		klog.Error("查询失败!", err)
		http.FailWithMessage("数据不存在", c)
	} else {
		http.OkWithData(property, c)
	}
}

// GetModulePropertyList 分页获取moduleProperty表列表
// @Tags ModuleProperty
// @Summary 分页获取moduleProperty表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query testReq.ModulePropertySearch true "分页获取moduleProperty表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /module/getModulePropertyList [get]
func GetModulePropertyList(c *gin.Context) {
	var pageInfo request.ModulePropertySearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := service.GetModulePropertyInfoList(&pageInfo); err != nil {
		klog.Error("获取失败!", err)
		http.FailWithMessage("获取失败", c)
	} else {
		if list == nil {
			list = &[]dao.ModuleProperty{}
		}
		http.OkWithDetailed(http.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
