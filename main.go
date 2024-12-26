package main

import (
	"heychat/router"
	"heychat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	r := router.Router()
	r.Run(":8081")
}
