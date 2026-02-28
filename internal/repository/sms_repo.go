package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/hacker4257/pet_charity/internal/database"
)

type SmsRepo struct{}

func NewSmsRepo() *SmsRepo {
	return &SmsRepo{}
}

// redis key 格式: sms:login:xxxx
func (r *SmsRepo) smsKey(purpose, phone string) string {
	return fmt.Sprintf("sms:%s:%s", purpose, phone)
}

//限流key， 60秒内不能重复复发
func (r *SmsRepo) throttleKey(phone string) string {
	return fmt.Sprintf("sms:throttle:%s", phone)
}

//每日发送次数
func (r *SmsRepo) dailyKey(phone string) string {
	return fmt.Sprintf("sms:daily:%s", phone)
}

//保存验证码到redis
func (r *SmsRepo) SaveCode(purpose, phone, code string) error {
	ctx := context.Background()
	key := r.smsKey(purpose, phone)
	return database.RDB.Set(ctx, key, code, 5*time.Minute).Err()
}

//获取验证码
func (r *SmsRepo) GetCode(purpose, phone string) (string, error) {
	ctx := context.Background()
	key := r.smsKey(purpose, phone)
	return database.RDB.Get(ctx, key).Result()
}

func (r *SmsRepo) DeleteCode(purpose, phone string) error {
	ctx := context.Background()
	key := r.smsKey(purpose, phone)
	return database.RDB.Del(ctx, key).Err()
}

//检查60s限流
func (r *SmsRepo) CheckThrottle(phone string) (bool, error) {
	ctx := context.Background()
	key := r.throttleKey(phone)
	exists, err := database.RDB.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

//设置60s限流
func (r *SmsRepo) SetThrottle(phone string) error {
	ctx := context.Background()
	key := r.throttleKey(phone)
	return database.RDB.Set(ctx, key, 1, 60*time.Second).Err()
}

//增加每次发送次数，返回当前次数
func (r *SmsRepo) IncrDaily(phone string) (int64, error) {
	ctx := context.Background()
	key := r.dailyKey(phone)

	count, err := database.RDB.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	//第一次设置时， 给key加上过期时间到当天结束
	if count == 1 {
		now := time.Now()
		endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		database.RDB.ExpireAt(ctx, key, endOfDay)
	}
	return count, nil
}
