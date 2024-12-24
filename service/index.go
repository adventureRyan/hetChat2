package service

import (
	"github.com/gin-gonic/gin"
)

// gin.Context 是 Gin 框架中的核心对象，包含了 HTTP 请求的上下文，例如请求信息、响应信息和中间件的数据。
// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "welcome!!",
	})
}
