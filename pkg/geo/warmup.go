package geo

import (
	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/hacker4257/pet_charity/pkg/logger"
)

func WarmupGeo() {
	//加载已审核通过的机构
	var orgs []model.Organization

	database.DB.Select("id, longitude, latitude").
		Where("status = ? AND  longitude != 0 AND latitude != 0", "approved").Find(&orgs)
	logger.Infof("[geo] warmed up %d organizations", len(orgs))

	for _, org := range orgs {
		Add(KeyOrgs, org.ID, org.Longitude, org.Latitude)
	}

	//加载未关闭的救助信息
	var rescues []model.Rescue
	database.DB.Select("id, longitude, latitude").
		Where("status != ? AND longitude != 0 AND latitude != 0", "closed").Find(&rescues)
	for _, r := range rescues {
		Add(KeyRescues, r.ID, r.Longitude, r.Latitude)
	}
	logger.Infof("[geo] warmed up %d rescues", len(rescues))
}
