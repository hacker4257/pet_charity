package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

const LeaderboardKey = "leaderboard:activity"

type ActivityRepo struct {
	db  *gorm.DB
	rdb *redis.Client
	sf  singleflight.Group
}

func NewActivityRepo() *ActivityRepo {
	return &ActivityRepo{
		db:  database.DB,
		rdb: database.RDB,
	}
}

// CreateLog 写入活跃度日志
func (r *ActivityRepo) CreateLog(log *model.UserActivityLogs) error {
	return r.db.Create(log).Error
}

// AddScore 给用户总分做原子加法（MySQL）
func (r *ActivityRepo) AddScore(userID uint, points int) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		UpdateColumn("activity_score", gorm.Expr("activity_score + ?", points)).Error
}

// RedisAddScore 同步加分到 Redis 排行榜
func (r *ActivityRepo) RedisAddScore(userID uint, points int) error {
	return r.rdb.ZIncrBy(
		context.Background(),
		LeaderboardKey,
		float64(points),
		strconv.Itoa(int(userID)),
	).Err()
}

// GetTopN 取排行榜前 N 名，返回 userID 和分数
func (r *ActivityRepo) GetTopN(page, pageSize int) ([]redis.Z, error) {
	start := int64((page - 1) * pageSize)
	stop := start + int64(pageSize) - 1
	return r.rdb.ZRevRangeWithScores(
		context.Background(),
		LeaderboardKey,
		start,
		stop,
	).Result()
}

// GetUserRank 查某用户排名（从1开始）和分数
func (r *ActivityRepo) GetUserRank(userID uint) (rank int64, score float64, err error) {
	ctx := context.Background()
	member := strconv.Itoa(int(userID))

	// 排名（ZRevRank 从0开始，+1 变成从1开始）
	rank, err = r.rdb.ZRevRank(ctx, LeaderboardKey, member).Result()
	if err != nil {
		return 0, 0, err
	}
	rank++

	// 分数
	score, err = r.rdb.ZScore(ctx, LeaderboardKey, member).Result()
	if err != nil {
		return 0, 0, err
	}

	return rank, score, nil
}

// GetUsersByIDs 批量查用户信息（排行榜展示用）
func (r *ActivityRepo) GetUsersByIDs(ids []uint) ([]model.User, error) {
	var users []model.User
	if len(ids) == 0 {
		return users, nil
	}
	ctx := context.Background()
	var missIDs []uint
	//1.先从redis
	for _, id := range ids {
		result, err := r.rdb.HGetAll(ctx, "user:cache:"+strconv.Itoa(int(id))).Result()
		if err != nil || len(result) == 0 {
			missIDs = append(missIDs, id)
			continue
		}
		if result["nickname"] == "" && result["avatar"] == "" {
			continue
		}
		users = append(users, model.User{
			BaseModel: model.BaseModel{ID: id},
			Nickname:  result["nickname"],
			Avatar:    result["avatar"],
		})
	}
	//2.未命中查MySQL
	if len(missIDs) > 0 {
		var dbUsers []model.User
		if err := r.db.Where("id IN ?", missIDs).Find(&dbUsers).Error; err != nil {
			return nil, err
		}
		//3回写缓存
		//把查到的用户做成map
		foundMap := make(map[uint]bool)
		for _, u := range dbUsers {
			foundMap[u.ID] = true
			key := "user:cache:" + strconv.Itoa(int(u.ID))
			r.rdb.HSet(ctx, key, "nickname", u.Nickname, "avatar", u.Avatar)
			r.rdb.Expire(ctx, key, 30*time.Minute)
		}
		//没有查到写空标记
		for _, id := range missIDs {
			if !foundMap[id] {
				key := "user:cache:" + strconv.Itoa(int(id))
				r.rdb.HSet(ctx, key, "nickname", "", "avatar", "")
				r.rdb.Expire(ctx, key, 5*time.Minute)
			}
		}
		users = append(users, dbUsers...)
	}

	return users, nil
}
