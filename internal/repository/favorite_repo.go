package repository

import (
	"context"
	"strconv"

	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/redis/go-redis/v9"
)

type FavoriteRepo struct {
	rdb *redis.Client
}

func NewFavoriteRepo() *FavoriteRepo {
	return &FavoriteRepo{rdb: database.RDB}
}

//petKey 某个宠物被哪些用户收藏
func petKey(petID uint) string {
	return "pet:fav:" + strconv.Itoa(int(petID))
}

//userKey 某个用户收藏了哪些宠物
func userKey(userID uint) string {
	return "user:favs:" + strconv.Itoa(int(userID))
}

//add 收藏
func (r *FavoriteRepo) Add(userID, petID uint) error {
	ctx := context.Background()
	pipe := r.rdb.Pipeline()
	pipe.SAdd(ctx, petKey(petID), userID)
	pipe.SAdd(ctx, userKey(userID), petID)
	_, err := pipe.Exec(ctx)
	return err
}

//取消收藏
func (r *FavoriteRepo) Remove(userID, petID uint) error {
	ctx := context.Background()
	pipe := r.rdb.Pipeline()
	pipe.SRem(ctx, petKey(petID), userID)
	pipe.SRem(ctx, userKey(userID), petID)
	_, err := pipe.Exec(ctx)
	return err
}

//是否已收藏
func (r *FavoriteRepo) IsFavorited(userID, petID uint) (bool, error) {
	return r.rdb.SIsMember(context.Background(), petKey(petID), userID).Result()
}

//某个宠物收藏总数
func (r *FavoriteRepo) PetFavCount(petID uint) (int64, error) {
	return r.rdb.SCard(context.Background(), petKey(petID)).Result()
}

//用户收藏的所有宠物ID
func (r *FavoriteRepo) UserFavPetIDs(userID uint) ([]uint, error) {
	members, err := r.rdb.SMembers(context.Background(), userKey(userID)).Result()
	if err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(members))
	for _, m := range members {
		id, _ := strconv.Atoi(m)
		ids = append(ids, uint(id))
	}
	return ids, nil
}
