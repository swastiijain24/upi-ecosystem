package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client{

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})


	_, err := rdb.Ping(context.Background()).Result()
	if err != nil{
		log.Fatal("redis connection failed:", err)
	}

	log.Println("redis connected")

	return rdb 
}

