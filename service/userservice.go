package service

import (
	"fmt"
	"heychat/models"
	"heychat/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/rand"
)

// GetUserList
// @Summary 用户列表
// @Tags 用户模块
// @Success 200 {string} json{"code", "message"}
// @Router /user/getUserList [post]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"message": data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param repassword formData string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.PostForm("name")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")
	salt := fmt.Sprintf("%06d", rand.Int31())
	if password != repassword {
		c.JSON(-1, gin.H{
			"message": "两次密码不一致",
		})
		return
	}
	// user.PassWord = password
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	if models.FindUserByName(user.Name).ID != 0 {
		c.JSON(400, gin.H{
			"message": "用户名已存在",
		})
		return
	}
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"message": "新增用户成功",
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id formData string false "用户名"
// @Success 200 {string} json{"code","message"}
// @Router /user/DeleteUser [post]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"message": "删除用户成功",
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param phone formData string false "电话"
// @param email formData string false "邮箱"
// @Success 200 {string} json{"code","message"}
// @Router /user/UpdateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	// _, err := govalidator.ValidateStruct(user)
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println("修改参数不匹配:", err)
		c.JSON(400, gin.H{
			"message": "修改参数不匹配",
		})
		return
	}
	models.UpdateUser(user)
	c.JSON(200, gin.H{
		"message": "修改用户成功",
	})
}

// UserLogin
// @Summary 登录
// @Tags 用户模块
// @param name formData string false "name"
// @param password formData string false "password"
// @Success 200 {string} json{"code","message"}
// @Router /user/UserLogin [post]
func UserLogin(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	u := models.FindUserByName(user.Name)
	if u.Name == "" {
		c.JSON(200, gin.H{
			"message": "没有这个用户",
		})
		return
	}
	salt := u.Salt
	if !utils.ValidPassword(user.PassWord, salt, u.PassWord) {
		c.JSON(200, gin.H{
			"message": "密码错误",
		})
		return
	}

	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.Md5Encode(str)
	utils.DB.Model(&u).Where("id=?", u.ID).Update("identity", temp)
	c.JSON(200, gin.H{
		"code":    1, // 1为成功,0为失败
		"message": "登录成功",
		"data":    u,
	})
}

// 防止跨域站点伪造请求
// 将 HTTP 升级为 WebSocket 的配置结构
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 表示允许所有来源进行连接
	},
}

// 通过 Gin 框架接收请求，将 HTTP 升级为 WebSocket，并处理消息逻辑。
func SendMsg(c *gin.Context) {
	// 升级 HTTP 到 WebSocket
	// c.Writer 和 c.Request 分别是 Gin 框架提供的 HTTP 响应写入器和请求。
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	fmt.Println("OK")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}

// 处理 WebSocket 连接的具体逻辑，将消息通过 WebSocket 推送到客户端。
func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	msg, err := utils.Subscribe(c, utils.PublishKey)
	if err != nil {
		fmt.Println(err)
	}
	tm := time.Now().Format("2006-01-02 15:04:05")
	m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
	// 发送消息:将消息发送到 WebSocket 客户端。
	err = ws.WriteMessage(1, []byte(m))
	fmt.Println("你好啊", m)
	if err != nil {
		fmt.Println(err)
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
