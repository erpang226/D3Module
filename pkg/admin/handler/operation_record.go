package handler

import (
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	"net/http"
)

func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO just log
		if c.Request.Method != http.MethodGet {
			klog.Info("get method pass")
		}
		c.Next()
	}
}
