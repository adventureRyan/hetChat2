package service

import (
	"fmt"
	"heychat/models"
	"net/http"
	"strconv"
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

// 用于渲染聊天界面
func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/main.html",
		"views/chat/userinfo.html",
		"views/chat/createcom.html")
	if err != nil {
		// 模板解析失败时记录错误并返回500状态码
		c.String(http.StatusInternalServerError, "Error parsing templates: %v", err)
		return
	}

	// if err != nil {
	// 	panic(err)
	// }
	// 获取用户输入
	userId, _ := strconv.Atoi(c.Query("userId"))
	if err != nil {
		// 用户ID解析失败
		c.String(http.StatusBadRequest, "Invalid userId: %v", err)
		return
	}
	// 获取用户身份验证信息
	token := c.Query("token")

	// 构造用户对象
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token

	// 执行模板渲染
	// 使用 Execute 方法将模板渲染为 HTML 输出到 c.Writer（HTTP 响应流）。
	// 渲染时会将用户对象 user 作为模板的数据源。
	err = ind.Execute(c.Writer, user)
	if err != nil {
		// 模板执行失败时记录错误并返回500状态码
		c.String(http.StatusInternalServerError, "Error executing template: %v", err)
	}
	// c.JSON(200, gin.H{
	// 	"message": "ToChat Page!!",
	// })
}
