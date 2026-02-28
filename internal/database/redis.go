package database

import (
	"context"
	"fmt"

	"github.com/hacker4257/pet_charity/internal/config"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() error {
	cfg := config.Global.Redis

	RDB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     20,
		MinIdleConns: 5,
		MaxRetries:   3,
	})

	//测试连接
	if err := RDB.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("Connect redis failed: %w", err)
	}

	return nil
}
