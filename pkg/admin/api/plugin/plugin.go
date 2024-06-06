package plugin

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	"main/common/types/http"
	"main/pkg/admin/api/plugin/request"
	service "main/pkg/admin/service/plugin"
	dao "main/pkg/common/dbm/dao/plugin"
)

// CreatePlugin 创建Plugin表
// @Tags Plugin
// @Summary 创建Plugin表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body Plugin true "创建Plugin表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /Plugin/createPlugin [post]
func CreatePlugin(c *gin.Context) {
	var m dao.Plugin
	err := c.ShouldBindJSON(&m)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.CreatePlugin(&m); err != nil {
		klog.Error("创建失败!", err)
		http.FailWithMessage(err.Error(), c)
	} else {
		http.OkWithData(int(m.Id), c)
	}
}

// DeletePlugin 删除Plugin表
// @Tags Plugin
// @Summary 删除Plugin表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body Plugin true "删除Plugin表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /Plugin/deletePlugin [delete]
func DeletePlugin(c *gin.Context) {
	var m dao.Plugin
	err := c.ShouldBindJSON(&m)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.DeletePluginById(m.Id); err != nil {
		klog.Error("删除失败!", err)
		http.FailWithMessage("删除失败", c)
	} else {
		http.OkWithMessage("删除成功", c)
	}
}

// DeletePluginByIds 批量删除Plugin表
// @Tags Plugin
// @Summary 批量删除Plugin表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Plugin表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /Plugin/deletePluginByIds [delete]
func DeletePluginByIds(c *gin.Context) {
	var IDS http.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.DeletePluginByIds(IDS); err != nil {
		klog.Error("批量删除失败!", err)
		http.FailWithMessage("批量删除失败", c)
	} else {
		http.OkWithMessage("批量删除成功", c)
	}
}

// UpdatePlugin 更新Plugin表
// @Tags Plugin
// @Summary 更新Plugin表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body Plugin true "更新Plugin表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /Plugin/updatePlugin [put]
func UpdatePlugin(c *gin.Context) {
	var m request.PluginUpdate
	err := c.ShouldBindJSON(&m)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.UpdatePlugin(m); err != nil {
		klog.Error("更新失败!", err)
		http.FailWithMessage("更新失败", c)
	} else {
		http.OkWithMessage("更新成功", c)
	}
}

// FindPlugin 用id查询Plugin表
// @Tags Plugin
// @Summary 用id查询Plugin表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query Plugin true "用id查询Plugin表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /Plugin/findPlugin [get]
func FindPlugin(c *gin.Context) {
	var m dao.Plugin
	err := c.ShouldBindQuery(&m)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if data, err := service.GetPlugin(m.Id); err != nil {
		klog.Error("查询失败!", err)
		http.FailWithMessage("数据不存在", c)
	} else {
		http.OkWithData(data, c)
	}
}

// GetPluginList 分页获取Plugin表列表
// @Tags Plugin
// @Summary 分页获取Plugin表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query testReq.PluginSearch true "分页获取Plugin表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /Plugin/getPluginList [get]
func GetPluginList(c *gin.Context) {
	info := &request.PluginSearch{}
	err := c.Bind(info)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := service.GetPluginInfoList(info); err != nil {
		klog.Error("获取失败!", err)
		http.FailWithMessage("获取失败", c)
	} else {
		if list == nil {
			list = &[]dao.Plugin{}
		}
		http.OkWithDetailed(http.PageResult{
			List:     list,
			Total:    total,
			Page:     info.Page,
			PageSize: info.PageSize,
		}, "获取成功", c)
	}
}
