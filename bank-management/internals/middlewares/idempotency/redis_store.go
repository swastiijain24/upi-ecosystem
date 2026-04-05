package idempotency

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client{

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil{
		log.Fatal("redis connection failed:", err)
	}

	log.Println("redis connected")

	return redisClient 
}

type RedisStore struct {
	redisClient *redis.Client
	ttl time.Duration
	prefix string
}

func NewRedisStore(redisClient *redis.Client, ttl time.Duration) *RedisStore{
	return &RedisStore{
		redisClient: redisClient,
		ttl: ttl,
		prefix: "idempotency:",
	}
}


func(s* RedisStore) key(idempotencyKey string) string{
	return s.prefix + idempotencyKey
}

func (s* RedisStore) Get(key string) (*Response, error){
	ctx:= context.Background()

	data, err := s.redisClient.Get(ctx, s.key(key)).Bytes()
	if errors.Is(err, redis.Nil){
		return nil , nil 
	}
	if err != nil{
		return nil, err
	}

	var resp Response 
	if err:= json.Unmarshal(data, &resp); err !=nil{
		return nil, err

	}

	return &resp, nil 
}

func (s* RedisStore) Set(key string, response *Response) error{
	ctx :=  context.Background()

	response.CreatedAt = time.Now()
	data, err := json.Marshal(response)
	if err!=nil{
		return err 
	}

	return s.redisClient.Set(ctx, s.key(key), data, s.ttl).Err()
}

func (s* RedisStore) SetProcessing(key string) (bool, error){
	ctx:= context.Background()

	ok, err := s.redisClient.SetNX(ctx, s.key(key), `{"status_code":0}`, s.ttl).Result()
	if err !=nil{
		return false, err
	}

	return ok, nil 
}

func (s* RedisStore) Delete(key string) error {
	ctx:= context.Background()
	return s.redisClient.Del(ctx, s.key(key)).Err()
}