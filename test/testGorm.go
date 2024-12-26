package main

import (
	"fmt"
	"heychat/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/heyChat2?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("fail to connect to the database")
	}
	// fmt.Println("Connect to the database")

	// 传递的参数应该是指针
	// 根据提供的模型结构体，自动创建或更新数据库中的表
	// db.AutoMigrate(&models.UserBasic{})
	db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.GroupBasic{})
	db.AutoMigrate(&models.Contact{})

	user := &models.UserBasic{
		Name:          "王五",
		LoginTime:     uint64(time.Now().Unix()),
		HeartbeatTime: uint64(time.Now().Unix()),
		LogOutTime:    uint64(time.Now().Unix()),
	}
	// fmt.Println(reflect.TypeOf(user))

	db.Create(user)

	db.Model(user).Update("PassWord", "1234")
	fmt.Println("创建用户表成功")
}
