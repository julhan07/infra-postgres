package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Connection struct {
	Host     string
	Port     int
	User     string
	Password string
}

var ctx = context.Background()

type Redis struct {
	Client *redis.Client
}

func NewRedis(conf Connection) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{Client: client}, nil
}

// SaveDataToRedis menyimpan data dalam format JSON di Redis
func (r *Redis) SaveDataToRedis(key string, data interface{}, duration time.Duration) error {

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	if err := r.Client.Set(ctx, key, jsonBytes, duration).Err(); err != nil {
		return err
	}

	return nil
}

// RetrieveDataFromRedis mengambil data dari Redis dan mengembalikannya dalam format JSON
func (r *Redis) RetrieveDataFromRedis(key string, result interface{}) error {
	data, err := r.Client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	return nil
}

// DeleteDataFromRedis menghapus data dari Redis berdasarkan kunci
func (r *Redis) DeleteDataFromRedis(key string) error {
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
