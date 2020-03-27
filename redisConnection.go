package comet

import (
	"fmt"

	"github.com/go-redis/redis"
)

var storage *redis.Client

func NewRedisStorage(url string, password string) *redis.Client {
	var options = redis.Options{
		Password: password,
		DB:       0,
		Addr:     url,
		OnConnect: func(conn *redis.Conn) error {
			go fmt.Println("[new redis connection]")
			return nil
		},
	}
	return redis.NewClient(&options)
}
