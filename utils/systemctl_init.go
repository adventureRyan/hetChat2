package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	// 自定义日志模板,打印SQL语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // 级别
			Colorful:      true,        // 彩色
		},
	)
	// dsn := "root:root@tcp(127.0.0.1:3306)/heyChat2?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := viper.GetString("mysql.dns")
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		fmt.Println("Init Mysql failed", err)
	}
}
