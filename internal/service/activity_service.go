package service

import (
	"errors"
	"strconv"

	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/hacker4257/pet_charity/internal/repository"
	"github.com/hacker4257/pet_charity/pkg/event"
	"github.com/hacker4257/pet_charity/pkg/logger"
)

var pointsRule = map[string]int{
	"login":    1,
	"donation": 10,
	"adoption": 5,
	"rescue":   3,
	"chat":     1,
}

type ActivityService struct {
	activityRepo repository.ActivityRepository
}

func NewActivityService(repo repository.ActivityRepository) *ActivityService {
	return &ActivityService{activityRepo: repo}
}

//记录活跃度并加分
func (s *ActivityService) AddActivity(userID uint, action string) error {
	points, ok := pointsRule[action]
	if !ok {
		return errors.New("unknown action")
	}

	//写日志
	log := &model.UserActivityLogs{
		UserID: userID,
		Action: action,
		Points: points,
	}
	if err := s.activityRepo.CreateLog(log); err != nil {
		return err
	}

	//2. MySQL 总分+n
	if err := s.activityRepo.AddScore(userID, points); err != nil {
		return err
	}

	//3.redis排行榜
	if err := s.activityRepo.RedisAddScore(userID, points); err != nil {
		return err
	}

	return nil
}

type LeaderboardItem struct {
	UserID   uint   `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Score    int    `json:"score"`
	Rank     int    `json:"rank"`
}

//获取排行榜
func (s *ActivityService) GetLeaderboard(page, pagesize int) ([]LeaderboardItem, error) {
	//1.从redis取topn
	results, err := s.activityRepo.GetTopN(page, pagesize)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return []LeaderboardItem{}, nil
	}

	//2.提取userID 列表
	userIDs := make([]uint, 0, len(results))
	for _, z := range results {
		id, _ := strconv.Atoi(z.Member.(string))
		userIDs = append(userIDs, uint(id))
	}

	//3.批量查询用户信息
	users, err := s.activityRepo.GetUsersByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	//4.做一个map
	userMap := make(map[uint]model.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	//5.组装结果
	startRank := (page-1)*pagesize + 1
	items := make([]LeaderboardItem, 0, len(results))
	for i, z := range results {
		id, _ := strconv.Atoi(z.Member.(string))
		user := userMap[uint(id)]
		items = append(items, LeaderboardItem{
			UserID:   uint(id),
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Score:    int(z.Score),
			Rank:     startRank + i,
		})

	}

	return items, nil
}

type MyRankInfo struct {
	Rank  int64   `json:"rank"`
	Score float64 `json:"score"`
}

func (s *ActivityService) GetMyRank(userID uint) (*MyRankInfo, error) {
	rank, score, err := s.activityRepo.GetUserRank(userID)
	if err != nil {
		return nil, err
	}
	return &MyRankInfo{Rank: rank, Score: score}, nil
}

//注册事件
func (s *ActivityService) RegisterHook() {
	event.Subscribe(func(e event.Event) {
		if _, ok := pointsRule[e.Action]; !ok {
			return
		}
		if err := s.AddActivity(e.UserID, e.Action); err != nil {
			logger.Error("handle event failed",
				logger.Str("action", e.Action),
				logger.Uint("user_id", e.UserID),
				logger.Err(err),
			)
		}
	})
}
