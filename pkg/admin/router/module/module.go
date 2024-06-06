package module

import (
	"github.com/gin-gonic/gin"
	"main/pkg/admin/api/module"
	"main/pkg/admin/api/moduleproperty"
	"main/pkg/admin/handler"
)

// InitRouter 初始化 module表 路由信息
func InitRouter(Router *gin.RouterGroup) {
	moduleRouter := Router.Group("module").Use(handler.OperationRecord())
	moduleRouterWithoutRecord := Router.Group("module")
	{
		moduleRouter.POST("createModule", module.CreateModule)             // 新建module表
		moduleRouter.DELETE("deleteModule", module.DeleteModule)           // 删除module表
		moduleRouter.DELETE("deleteModuleByIds", module.DeleteModuleByIds) // 批量删除module表
		moduleRouter.PUT("updateModule", module.UpdateModule)              // 更新module表
		moduleRouter.POST("startModule", module.StartModule)               // 启动模块
	}
	{
		moduleRouterWithoutRecord.GET("findModule", module.FindModule)                               // 根据ID获取module表
		moduleRouterWithoutRecord.GET("findModuleAndProperty", module.FindModuleAndProperty)         // 根据ID获取module表
		moduleRouterWithoutRecord.GET("getModuleListWithProperty", module.GetModuleListWithProperty) // 根据ID获取module表
		moduleRouterWithoutRecord.GET("findModuleByName", module.FindModuleByName)                   // 根据name获取module表
		moduleRouterWithoutRecord.GET("getModuleList", module.GetModuleList)                         // 获取module表列表
		moduleRouterWithoutRecord.GET("getModuleListRouter", module.GetModuleListRouter)             // 获取module表列表
	}
	{
		moduleRouter.POST("createModuleProperty", moduleproperty.CreateModuleProperty)             // 新建module property表
		moduleRouter.DELETE("deleteModuleProperty", moduleproperty.DeleteModuleProperty)           // 删除module property表
		moduleRouter.DELETE("deleteModulePropertyByIds", moduleproperty.DeleteModulePropertyByIds) // 批量删除module property表
		moduleRouter.PUT("updateModuleProperty", moduleproperty.UpdateModuleProperty)              // 更新module property表
	}
	{
		moduleRouterWithoutRecord.GET("findModuleProperty", moduleproperty.FindModuleProperty)       // 根据ID获取module property表
		moduleRouterWithoutRecord.GET("getModulePropertyList", moduleproperty.GetModulePropertyList) // 获取module property表列表
	}
}
