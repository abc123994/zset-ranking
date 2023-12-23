package utils

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

type MemberScore struct {
	Member string
	Score  float64
}

func InitRedis() {
	pool = &redis.Pool{
		MaxIdle:   5,
		MaxActive: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}
func GetAllScoresWithMembers(zsetKey string) ([]MemberScore, error) {
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
