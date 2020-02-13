package utils

import (
	"github.com/garyburd/redigo/redis"
)

var (
	Redis redis.Conn
)

func init() {
	conn, err := redis.Dial("tcp", ":6370")
	if err != nil {
		panic(err)
	}
	Redis = conn
}
