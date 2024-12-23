package main

import "heychat/router"

func main() {
	r := router.Router()
	r.Run(":8081")
}
