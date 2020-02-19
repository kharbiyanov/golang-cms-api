package utils

import (
	"github.com/garyburd/redigo/redis"
)

var (
	Redis redis.Conn
)

func init() {
	conn, err := redis.Dial("tcp", Config.RedisAddr)
	if err != nil {
		panic(err)
	}
	Redis = conn
}
