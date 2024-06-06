package handler

import (
	"github.com/gin-gonic/gin"
	"main/common/types/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 Authorization 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中
		// 不过需要协商过期时间 可以约定刷新令牌或者重新登录
		token, _ := c.Cookie(http.AuthorizationKey)
		if len(token) == 0 {
			token = c.Request.Header.Get(http.AuthorizationKey)
		}
		if token != "dev" {
			http.FailWithDetailed(gin.H{"reload": true}, "未登录或非法访问", c)
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Next()
	}
}
