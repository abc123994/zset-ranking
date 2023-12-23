package main

import (
	"log"
	"zorder/utils"
)

func main() {

	utils.InitRedis()

	data := utils.CacheLeaderBoards("leaderboards", "myzset", 300)
	log.Println(data)
}
