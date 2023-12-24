package utils

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func InitRedis() {
	pool = &redis.Pool{
		MaxIdle:   5,
		MaxActive: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}

type Ranking struct {
	Rank []Data            `json:"leaderboard"`
	Key  map[string]Member `json:"customer"`
}
type Data struct {
	Rank   int32  `json:"rank"`
	Score  int32  `json:"score"`
	Member string `json:"member"`
}
type Member struct {
	Score int32 `json:"score"`
	Rank  int32 `json:"rank"`
}

/*
		{
		  leaderboard:[{0} {1} {2}]        //get the board
		  customer: "a":{},"b":{},"c":{}   //mapping for query self

	    }
*/
func CacheLeaderBoards(key string, zkey string) (rank Ranking) {

	data, err := getLeaderBoards(key)

	if err == nil && data != "" {
		json.Unmarshal([]byte(data), &rank)
	} else {

		sc, _ := getAllScoresWithMembers(zkey)
		do_rank := Ranking{}
		do_rank.Key = map[string]Member{}
		_current := 0
		_score := 1

		for _, e := range sc {
			if e.Score != float64(_score) {
				_score = int(e.Score)
				_current++
				_tmp := Data{Rank: int32(_current), Score: int32(e.Score), Member: e.Member}
				do_rank.Key[e.Member] = Member{Score: int32(e.Score), Rank: int32(_current)}
				do_rank.Rank = append(do_rank.Rank, _tmp)
			} else {
				_tmp := Data{Rank: int32(_current), Score: int32(e.Score), Member: e.Member}
				do_rank.Key[e.Member] = Member{Score: int32(e.Score), Rank: int32(_current)}
				do_rank.Rank = append(do_rank.Rank, _tmp)
			}

		}
		b, err := json.Marshal(&do_rank)

		if err != nil {
			fmt.Println(err)

		}

		setLeaderBoards(key, string(b), 300)
		rank = do_rank
	}
	return
}
func setLeaderBoards(key string, val string, ttl int32) {

	conn := pool.Get()
	defer conn.Close()
	conn.Do("SET", key, val)
	conn.Do("ExPire", key, ttl)
}
func getLeaderBoards(key string) (string, error) {

	conn := pool.Get()
	defer conn.Close()
	data, err := conn.Do("GET", key)

	out, err := redis.String(data, err)

	return out, err
}

type MemberScore struct {
	Member string
	Score  float64
}

func getAllScoresWithMembers(zsetKey string) ([]MemberScore, error) {
	conn := pool.Get()
	defer conn.Close()

	// Retrieve all members and their scores in descending order
	scoresWithMembers, err := redis.Values(conn.Do("ZREVRANGE", zsetKey, 0, -1, "WITHSCORES"))
	if err != nil {
		return nil, err
	}

	var result []MemberScore
	for i := 0; i < len(scoresWithMembers); i += 2 {
		member, _ := redis.String(scoresWithMembers[i], nil)
		score, _ := redis.Float64(scoresWithMembers[i+1], nil)
		result = append(result, MemberScore{Member: member, Score: score})
	}

	return result, nil
}
func IncrScoreByMember(zsetKey string, score int32, member string) {
	conn := pool.Get()
	defer conn.Close()

	// Increment the score of the specified member
	_, err := conn.Do("ZINCRBY", zsetKey, float64(score), member)
	if err != nil {
		fmt.Println("Error incrementing score:", err)
	}
}
