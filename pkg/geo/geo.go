package geo

import (
	"context"
	"strconv"

	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/redis/go-redis/v9"
)

const (
	KeyOrgs    = "geo:orgs"
	KeyRescues = "geo:rescues"
)

type NearbyItem struct {
	ID   uint
	Dist float64
}

// add 在geo 添加一个坐标
func Add(key string, id uint, lng, lat float64) error {
	if lng == 0 && lat == 0 {
		return nil
	}
	return database.RDB.GeoAdd(context.Background(), key, &redis.GeoLocation{
		Name:      strconv.FormatUint(uint64(id), 10),
		Longitude: lng,
		Latitude:  lat,
	}).Err()
}

//删除
func Remove(key string, id uint) error {
	return database.RDB.ZRem(context.Background(), key, strconv.FormatUint(uint64(id), 10)).Err()
}

//search 按照中心点+半径搜素，返回id列表+距离， 按照升序排序
func Search(key string, lng, lat, redisuKm float64, count int) ([]NearbyItem, error) {
	results, err := database.RDB.GeoSearchLocation(context.Background(), key,
		&redis.GeoSearchLocationQuery{
			GeoSearchQuery: redis.GeoSearchQuery{
				Longitude:  lng,
				Latitude:   lat,
				Radius:     redisuKm,
				RadiusUnit: "km",
				Sort:       "ASC",
				Count:      count,
			},
		},
	).Result()
	if err != nil {
		return nil, err
	}
	items := make([]NearbyItem, 0, len(results))
	for _, loc := range results {
		id, err := strconv.ParseUint(loc.Name, 10, 64)
		if err != nil {
			continue
		}
		items = append(items, NearbyItem{
			ID:   uint(id),
			Dist: loc.Dist,
		})

	}
	return items, nil
}
