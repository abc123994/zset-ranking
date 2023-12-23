package main

import (
	"log"
	"zorder/utils"
)

func main() {

	utils.InitRedis()

	data := utils.CacheLeaderBoards("leaderboards", "myzset") // leaderboards 和 zset 會因為活動改變
	log.Println(data)
}
