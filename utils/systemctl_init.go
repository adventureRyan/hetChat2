package utils

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 让外部可以访问
var DB *gorm.DB

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Read config error")
	}
	fmt.Println("config.app:", viper.Get("app"))
	fmt.Println("config.mysql:", viper.Get("mysql"))
}

func InitMySQL() {
	// dsn := "root:root@tcp(127.0.0.1:3306)/heyChat2?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := viper.GetString("mysql.dns")
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Init Mysql failed", err)
	}
}
