package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/redis/go-redis/v9"
)

const FeedKey = "feed:global"
const FeedMaxLen = 20

type FeedItem struct {
	Action    string `json:"action"`
	UserID    uint   `json:"user_id"`
	Timestamp int64  `json:"timestamp"`
}

type FeedRepo struct {
	rdb *redis.Client
}

func NewFeedRepo() *FeedRepo {
	return &FeedRepo{
		rdb: database.RDB,
	}
}

// Push 往全站动态列表头部插入一条
func (r *FeedRepo) Push(action string, userID uint) error {
	ctx := context.Background()
	item := FeedItem{
		Action:    action,
		UserID:    userID,
		Timestamp: time.Now().Unix(),
	}
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	pipe := r.rdb.Pipeline()
	pipe.LPush(ctx, FeedKey, data)            //头部插入
	pipe.LTrim(ctx, FeedKey, 0, FeedMaxLen-1) // 只保留最新200条
	_, err = pipe.Exec(ctx)
	return err
}

// List 取最新 N 条动态
func (r *FeedRepo) List(page, pageSize int) ([]FeedItem, error) {
	start := int64((page - 1) * pageSize)
	stop := start + int64(pageSize) - 1
	results, err := r.rdb.LRange(context.Background(), FeedKey, start, stop).Result()
	if err != nil {
		return nil, err
	}
	items := make([]FeedItem, 0, len(results))
	for _, raw := range results {
		var item FeedItem
		if err := json.Unmarshal([]byte(raw), &item); err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}
