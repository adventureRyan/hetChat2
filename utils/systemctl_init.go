package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 让外部可以访问
var (
	DB  *gorm.DB
	Red *redis.Client
)

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

func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.PoolSize"),
		MinIdleConns: viper.GetInt("redis.MinIdleConns"),
	})

	// 测试redis是否正常工作
	// Background返回一个非空的Context。 它永远不会被取消，没有值，也没有期限。
	// 它通常在main函数，初始化和测试时使用，并用作传入请求的顶级上下文。
	ctx := context.Background()
	pong, err := Red.Ping(ctx).Result()
	if err != nil {
		fmt.Println("init redis .", err)
	} else {
		fmt.Println("init redis success", pong)
	}
}

const (
	PublishKey = "websocket" // 可以用作 Redis 中的频道名称或者关键标识。
)

// Publish发布消息到Redis
// channel 是Redis的频道名称
func Publish(ctx context.Context, channel string, msg string) error {
	err := Red.Publish(ctx, channel, msg).Err()
	return err
}

// Subscribe订阅消息到Redis
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Red.Subscribe(ctx, channel) // 创建一个订阅对象
	fmt.Println("Subscribe", ctx)
	msg, err := sub.ReceiveMessage(ctx) // 阻塞等待频道中发布的消息
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Subscribe success:", msg.Payload)
	return msg.Payload, err
}
