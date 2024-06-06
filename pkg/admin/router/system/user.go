package system

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	http2 "main/common/types/http"
	"main/pkg/admin/api/system"
	"net/http"
)

func InitRouter(Router *gin.RouterGroup) {
	group := Router.Group("/")
	{
		group.POST("/login", func(c *gin.Context) {
			user := &system.Login{}
			err := c.Bind(user)
			if err != nil {
				klog.Fatal("login error ", err)
				c.JSON(http.StatusBadRequest, "{}")
			}
			if user != nil && user.Username == "admin" && user.Password == "123456" {
				http2.OkWithDetailed("dev", "登录成功", c)
			}
		})
	}
}
