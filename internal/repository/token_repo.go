package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hacker4257/pet_charity/internal/database"
)

type TokenRepo struct{}

func NewTokenRepo() *TokenRepo {
	return &TokenRepo{}
}

func tokenVerKey(userID uint) string {
	return fmt.Sprintf("token_ver:%d", userID)
}

//获取当前版本号，不存在返回0
func (r *TokenRepo) GetVersion(userID uint) int {
	val, err := database.RDB.Get(context.Background(), tokenVerKey(userID)).Result()
	if err != nil {
		return 0
	}
	v, _ := strconv.Atoi(val)
	return v
}

//版本号+1，返回新的
func (r *TokenRepo) IncrVersion(userID uint) (int, error) {
	val, err := database.RDB.Incr(context.Background(), tokenVerKey(userID)).Result()
	if err != nil {
		return 0, err
	}
	return int(val), nil
}