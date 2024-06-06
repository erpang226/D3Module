package module

import (
	"github.com/gin-gonic/gin"
	core "main/pkg/actor/core"
	"net/http"
)

func InitTestRouter(Router *gin.RouterGroup) {
	moduleGroup := Router.Group("/test")
	{
		moduleGroup.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, core.GetModules())
		})
		moduleGroup.GET("/:name", func(c *gin.Context) {
			name := c.Param("name")
			moduleInfo := core.GetModules()[name]
			c.JSON(http.StatusOK, moduleInfo.GetModule())
		})
		moduleGroup.POST("/mqtt/:name/start", func(c *gin.Context) {
			name := c.Param("name")
			moduleMap := core.GetModules()
			moduleInfo := moduleMap[name]
			moduleInfo.GetModule().SetEnable(true)

			core.StartModule(name, moduleMap[name])
		})
		moduleGroup.DELETE("/mqtt/:name", func(c *gin.Context) {
			name := c.Param("name")
			moduleMap := core.GetModules()
			moduleInfo := moduleMap[name]
			moduleInfo.GetModule().Stop()
			c.JSON(http.StatusOK, "{}")
		})

	}
}
