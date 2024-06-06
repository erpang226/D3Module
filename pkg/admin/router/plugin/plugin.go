package plugin

import (
	"github.com/gin-gonic/gin"
	"main/pkg/admin/api/plugin"
	"main/pkg/admin/handler"
)

// InitRouter 初始化 plugin表 路由信息
func InitRouter(Router *gin.RouterGroup) {
	pluginRouter := Router.Group("plugin").Use(handler.OperationRecord())
	pluginRouterWithoutRecord := Router.Group("plugin")
	{
		pluginRouter.POST("createPlugin", plugin.CreatePlugin)             // 新建plugin表
		pluginRouter.DELETE("deletePlugin", plugin.DeletePlugin)           // 删除plugin表
		pluginRouter.DELETE("deletePluginByIds", plugin.DeletePluginByIds) // 批量删除plugin表
		pluginRouter.PUT("updatePlugin", plugin.UpdatePlugin)              // 更新plugin表
	}
	{
		pluginRouterWithoutRecord.GET("findPlugin", plugin.FindPlugin)       // 根据ID获取plugin表
		pluginRouterWithoutRecord.GET("getPluginList", plugin.GetPluginList) // 获取plugin表列表
	}
}
