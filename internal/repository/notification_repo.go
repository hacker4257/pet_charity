package repository

import (
	"context"
	"strconv"

	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type NotificationRepo struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewNotificationRepo() *NotificationRepo {
	return &NotificationRepo{
		db:  database.DB,
		rdb: database.RDB,
	}
}

//
func (r *NotificationRepo) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepo) ListByUser(userID uint, page, pageSize int) ([]model.Notification, int64, error) {
	var list []model.Notification
	var total int64

	db := r.db.Model(&model.Notification{}).Where("user_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *NotificationRepo) MarkRead(userID, id uint) error {
	return r.db.Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func (r *NotificationRepo) MarkAllRead(userID uint) error {
	return r.db.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

//redis 统计未读
func unreadKey(userID uint) string {
	return "notify:unread:" + strconv.Itoa(int(userID))
}

func (r *NotificationRepo) GetUnreadCount(userID uint) (int64, error) {
	val, err := r.rdb.Get(context.Background(), unreadKey(userID)).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, nil
}

func (r *NotificationRepo) ResetUnread(userID uint) error {
	return r.rdb.Set(context.Background(), unreadKey(userID), 0, 0).Err()
}

func (r *NotificationRepo) IncrUnread(userID uint) error {
	return r.rdb.Incr(context.Background(), unreadKey(userID)).Err()
}

func (r *NotificationRepo) DecrUnread(userID uint) error {
	ctx := context.Background()
	val, err := r.rdb.Decr(ctx, unreadKey(userID)).Result()
	if err != nil {
		return err
	}
	if val < 0 {
		r.rdb.Set(ctx, unreadKey(userID), 0, 0)
	}
	return nil
}

const NotifyChannel = "notify:push"

func (r *NotificationRepo) Publish(userID uint, data []byte) error {
	//发布到频道
	channel := NotifyChannel + ":" + strconv.Itoa(int(userID))
	return r.rdb.Publish(context.Background(), channel, data).Err()
}
