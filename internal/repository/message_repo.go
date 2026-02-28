package repository

import (
	"fmt"

	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/hacker4257/pet_charity/internal/model"
	"gorm.io/gorm"
)

type messageRepo struct {
	db *gorm.DB
}

func NewMessageRepo() MessageRepository {
	return &messageRepo{db: database.DB}
}

// 保存消息
func (r *messageRepo) Create(msg *model.Message) error {
	return r.db.Create(msg).Error
}

// 私聊历史（双向）
func (r *messageRepo) ListPrivate(userA, userB uint, page, pageSize int) ([]model.Message,
	int64, error) {
	var msgs []model.Message
	var total int64

	query := r.db.Model(&model.Message{}).
		Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id =?)",
			userA, userB, userB, userA).
		Where("room_id = ''")

	query.Count(&total)
	err := query.Preload("FromUser").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&msgs).Error

	return msgs, total, err
}

// 房间消息历史
func (r *messageRepo) ListByRoom(roomID string, page, pageSize int) ([]model.Message,
	int64, error) {
	var msgs []model.Message
	var total int64

	query := r.db.Model(&model.Message{}).Where("room_id = ?", roomID)
	query.Count(&total)
	err := query.Preload("FromUser").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&msgs).Error

	return msgs, total, err
}

// 会话列表（最近聊过的人）
func (r *messageRepo) ListConversations(userID uint) ([]model.Message, error) {
	var msgs []model.Message
	// 每个会话取最新一条
	subQuery := r.db.Model(&model.Message{}).
		Select("MAX(id) as id").
		Where("room_id = '' AND (from_user_id = ? OR to_user_id = ?)", userID, userID).
		Group(fmt.Sprintf("CASE WHEN from_user_id = %d THEN to_user_id ELSE from_user_id END", userID))

	err := r.db.Preload("FromUser").
		Where("id IN (?)", subQuery).
		Order("created_at DESC").
		Find(&msgs).Error

	return msgs, err
}
