package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/redis/go-redis/v9"
)

type CacheRepo struct {
	rdb *redis.Client
}

func NewCacheRepo() *CacheRepo {
	return &CacheRepo{
		rdb: database.RDB,
	}
}

func userCacheKey(userID uint) string {
	return "user:cache:" + strconv.Itoa(int(userID))
}

type CachedUser struct {
	Nickname string
	Avatar   string
}

//写入用户缓存，过期时间30分钟
func (r *CacheRepo) SetUser(userID uint, nickname, avatar string) error {
	ctx := context.Background()
	key := userCacheKey(userID)
	pipe := r.rdb.Pipeline()
	pipe.HSet(ctx, key, "nickname", nickname, "avatar", avatar)
	pipe.Expire(ctx, key, 30*time.Minute)
	_, err := pipe.Exec(ctx)
	return err
}

//读取缓存
func (r *CacheRepo) GetUser(userID uint) (*CachedUser, error) {
	result, err := r.rdb.HGetAll(context.Background(), userCacheKey(userID)).Result()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return &CachedUser{
		Nickname: result["nickname"],
		Avatar:   result["avatar"],
	}, nil
}

//删除用户缓存
func (r *CacheRepo) DeleteUser(userID uint) error {
	return r.rdb.Del(context.Background(), userCacheKey(userID)).Err()
}
