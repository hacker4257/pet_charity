package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const OrderExpireKey = "order:expire"

type DonationRepo struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewDonationRepo() *DonationRepo {
	return &DonationRepo{
		db:  database.DB,
		rdb: database.RDB,
	}
}

func (r *DonationRepo) Create(donation *model.Donation) error {
	return r.db.Create(donation).Error
}

func (r *DonationRepo) FindByID(id uint) (*model.Donation, error) {
	var donation model.Donation
	result := r.db.Preload("User").First(&donation, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &donation, nil
}

// 根据订单号查找
func (r *DonationRepo) FindByTradeNo(tradeNo string) (*model.Donation, error) {
	var donation model.Donation
	result := r.db.Where("trade_no = ?", tradeNo).First(&donation)
	if result.Error != nil {
		return nil, result.Error
	}
	return &donation, nil
}

// 更新指定字段
func (r *DonationRepo) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&model.Donation{}).Where("id = ?", id).Updates(fields).Error
}

// 用户捐赠记录
func (r *DonationRepo) ListByUser(userID uint, page, pageSize int) ([]model.Donation, int64, error) {
	var donations []model.Donation
	var total int64

	db := r.db.Model(&model.Donation{}).
		Where("user_id = ? AND payment_status = ?", userID, "paid")
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&donations).Error; err != nil {
		return nil, 0, err
	}

	return donations, total, nil
}

// 公开捐赠公示 数据脱敏
func (r *DonationRepo) ListPublic(targetType string, targetID uint, page, pageSize int) ([]model.Donation, int64, error) {
	var donations []model.Donation
	var total int64

	db := r.db.Model(&model.Donation{}).Where("payment_status = ?", "paid")
	if targetType != "" {
		db = db.Where("target_type = ?", targetType)
	}
	if targetID > 0 {
		db = db.Where("target_id = ?", targetID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := db.Preload("User").Offset(offset).Limit(pageSize).Order("id DESC").Find(&donations).Error; err != nil {
		return nil, 0, err
	}

	return donations, total, nil
}

// 捐赠统计
func (r *DonationRepo) Stats() (*DonationStats, error) {
	var stats DonationStats

	//总捐赠金额和笔数
	r.db.Model(&model.Donation{}).
		Where("payment_status = ?", "paid").
		Select("COALESCE(SUM(amount), 0) as total_amount, COUNT(*) as total_count").
		Scan(&stats)

	//今日捐赠
	today := time.Now().Format("2006-01-02")
	r.db.Model(&model.Donation{}).
		Where("payment_status = ? AND DATE(paid_at) = ?", "paid", today).
		Select("COALESCE(SUM(amount), 0) as today_amount, COUNT(*) as today_count").
		Scan(&stats)

	return &stats, nil
}

type DonationStats struct {
	TotalAmount int64 `json:"total_amount"`
	TotalCount  int64 `json:"total_count"`
	TodayAmount int64 `json:"today_amount"`
	TodayCount  int64 `json:"today_count"`
}

// UpdateFieldsWithStatus 带条件的更新，返回受影响行数
func (r *DonationRepo) UpdateFieldsWithStatus(id uint, expectStatus string, fields map[string]interface{}) (int64, error) {
	result := r.db.Model(&model.Donation{}).
		Where("id = ? AND payment_status = ?", id, expectStatus).
		Updates(fields)
	return result.RowsAffected, result.Error
}

// FindPendingByUser 查找用户对同一目标的未支付订单
func (r *DonationRepo) FindPendingByUser(userID uint, targetType string, targetID uint) (*model.Donation, error) {
	var donation model.Donation
	result := r.db.Where(
		"user_id = ? AND target_type = ? AND target_id = ? AND payment_status = ?",
		userID, targetType, targetID, "pending",
	).First(&donation)
	if result.Error != nil {
		return nil, result.Error
	}
	return &donation, nil
}

// AddExpireTask 订单创建时，加入过期队列
func (r *DonationRepo) AddExpireTask(tradeNo string, expireAt time.Time) error {
	return r.rdb.ZAdd(context.Background(), OrderExpireKey, redis.Z{
		Score:  float64(expireAt.Unix()),
		Member: tradeNo,
	}).Err()
}

// GetExpiredTasks 取出所有已到期的订单号
func (r *DonationRepo) GetExpiredTasks() ([]string, error) {
	now := float64(time.Now().Unix())
	return r.rdb.ZRangeArgs(context.Background(), redis.ZRangeArgs{
		Key:     OrderExpireKey,
		Start:   "-inf",
		Stop:    strconv.FormatFloat(now, 'f', 0, 64),
		ByScore: true,
	}).Result()
}

// RemoveExpireTask 处理完后移除
func (r *DonationRepo) RemoveExpireTask(tradeNo string) error {
	return r.rdb.ZRem(context.Background(), OrderExpireKey, tradeNo).Err()
}
