package module

import (
	"main/common/types/http"
	"main/pkg/admin/api/module/request"
	"main/pkg/admin/api/module/response"
	service "main/pkg/admin/service/module"
	dao "main/pkg/common/dbm/dao/module"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

func CreateModule(c *gin.Context) {
	var m request.ModuleCreate
	err := c.ShouldBindJSON(&m)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.CreateModule(&m); err != nil {
		klog.Error("创建失败!", err)
		http.FailWithMessage(err.Error(), c)
	} else {
		http.OkWithData(int(m.Id), c)
	}
}

func DeleteModule(c *gin.Context) {
	var m dao.Module
	err := c.ShouldBindJSON(&m)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.DeleteModule(m.Id); err != nil {
		klog.Error("删除失败!", err)
		http.FailWithMessage("删除失败", c)
	} else {
		http.OkWithMessage("删除成功", c)
	}
}

func DeleteModuleByIds(c *gin.Context) {
	var IDS http.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.DeleteModuleByIds(IDS); err != nil {
		klog.Error("批量删除失败!", err)
		http.FailWithMessage("批量删除失败", c)
	} else {
		http.OkWithMessage("批量删除成功", c)
	}
}

func UpdateModule(c *gin.Context) {
	var m request.ModuleUpdate
	err := c.ShouldBindJSON(&m)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.UpdateModule(m); err != nil {
		klog.Error("更新失败!", err)
		http.FailWithMessage("更新失败", c)
	} else {
		http.OkWithMessage("更新成功", c)
	}
}

func FindModule(c *gin.Context) {
	var m dao.Module
	err := c.ShouldBindQuery(&m)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if data, err := service.GetModule(m.Id); err != nil {
		klog.Error("查询失败!", err)
		http.FailWithMessage("数据不存在", c)
	} else {
		http.OkWithData(data, c)
	}
}

func FindModuleAndProperty(c *gin.Context) {
	var m dao.Module
	err := c.ShouldBindQuery(&m)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if data, err := service.GetModuleAndProperty(m.Id); err != nil {
		klog.Error("查询失败!", err)
		http.FailWithMessage("数据不存在", c)
	} else {
		http.OkWithData(data, c)
	}
}

func FindModuleByName(c *gin.Context) {
	var m dao.Module
	err := c.ShouldBindQuery(&m)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if data, err := service.GetModuleByName(m.Name); err != nil {
		klog.Error("查询失败,或者数据不存在!", err)
		http.OkWithData(data, c)
	} else {
		http.OkWithData(data, c)
	}
}

func GetModuleList(c *gin.Context) {
	info := &request.ModuleSearch{}
	err := c.Bind(info)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := service.GetModuleInfoList(info); err != nil {
		klog.Error("获取失败!", err)
		http.FailWithMessage("获取失败", c)
	} else {
		if list == nil {
			list = &[]dao.Module{}
		}
		http.OkWithDetailed(http.PageResult{
			List:     list,
			Total:    total,
			Page:     info.Page,
			PageSize: info.PageSize,
		}, "获取成功", c)
	}
}

func GetModuleListWithProperty(c *gin.Context) {
	info := &request.ModuleSearch{}
	err := c.Bind(info)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := service.GetModuleListWithProperty(info); err != nil {
		klog.Error("获取失败!", err)
		http.FailWithMessage("获取失败", c)
	} else {
		if *list == nil {
			list = &[]response.DmpGatewaySubDeviceModuleResponse{}
		}
		http.OkWithDetailed(http.PageResult{
			List:     list,
			Total:    total,
			Page:     info.Page,
			PageSize: info.PageSize,
		}, "获取成功", c)
	}
}

func GetModuleListRouter(c *gin.Context) {
	info := &request.ModuleSearch{}
	err := c.Bind(info)
	if err != nil {
		http.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := service.GetModuleInfoListRouter(info); err != nil {
		klog.Error("获取失败!", err)
		http.FailWithMessage("获取失败", c)
	} else {
		if list == nil {
			list = &[]response.ModuleRouterResponse{}
		}
		http.OkWithDetailed(http.PageResult{
			List:     list,
			Total:    total,
			Page:     info.Page,
			PageSize: info.PageSize,
		}, "获取成功", c)
	}
}

func StartModule(c *gin.Context) {

}
