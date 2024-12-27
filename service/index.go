package service

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

// gin.Context 是 Gin 框架中的核心对象，包含了 HTTP 请求的上下文，例如请求信息、响应信息和中间件的数据。
// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index1.html")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading template: %v", err))
		return
	}
	err = ind.Execute(c.Writer, "index1")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error executing template: %v", err))
	}
	// c.JSON(200, gin.H{
	// 	"message": "welcome!!",
	// })
}

func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error loading template: %v", err))
		return
	}
	err = ind.Execute(c.Writer, "register")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error executing template: %v", err))
	}
	// c.JSON(200, gin.H{
	// 	"message": "welcome!!",
	// })
}
