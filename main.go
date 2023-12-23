package main

import (
	"log"
	"zorder/utils"
)

type ranking struct {
	rank []data
	k    map[string]member
}
type data struct {
	rank   int32
	score  int32
	member string
}
type member struct {
	score int32
	rank  int32
}

func main() {

	utils.InitRedis()

	sc, _ := utils.GetAllScoresWithMembers("myzset")
	do_rank := ranking{}
	do_rank.k = map[string]member{}
	_current := 0
	_score := 1

	for _, e := range sc {
		if e.Score != float64(_score) {
			_score = int(e.Score)
			_current++
			_tmp := data{rank: int32(_current), score: int32(e.Score), member: e.Member}
			do_rank.k[e.Member] = member{score: int32(e.Score), rank: int32(_current)}
			do_rank.rank = append(do_rank.rank, _tmp)
		} else {
			_tmp := data{rank: int32(_current), score: int32(e.Score), member: e.Member}
			do_rank.k[e.Member] = member{score: int32(e.Score), rank: int32(_current)}
			do_rank.rank = append(do_rank.rank, _tmp)
		}

	}
	log.Println(do_rank)

}
