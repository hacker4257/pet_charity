package service

import (
	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/hacker4257/pet_charity/internal/repository"
)

type FavoriteService struct {
	favRepo repository.FavoriteRepository
	petRepo repository.PetRepository
}

func NewFavoritService(favRepo repository.FavoriteRepository, petRepo repository.PetRepository) *FavoriteService {
	return &FavoriteService{
		favRepo: favRepo,
		petRepo: petRepo,
	}
}

//取消/收藏
func (s FavoriteService) Toggle(userID, petID uint) (bool, error) {
	//先检查是否已收藏
	faved, err := s.favRepo.IsFavorited(userID, petID)
	if err != nil {
		return false, err
	}

	if faved {
		//已收藏 -- 取消
		err = s.favRepo.Remove(userID, petID)
		return false, err

	}
	err = s.favRepo.Add(userID, petID)
	return true, nil
}

//收藏返回状态结构
type FavStatus struct {
	Favorited bool  `json:"favorited"`
	Count     int64 `json:"count"`
}

//查看宠物的收藏状态
func (s *FavoriteService) GetStatus(userID, petID uint) (*FavStatus, error) {
	faved, err := s.favRepo.IsFavorited(userID, petID)
	if err != nil {
		return nil, err
	}
	count, err := s.favRepo.PetFavCount(petID)
	if err != nil {
		return nil, err
	}
	return &FavStatus{Favorited: faved, Count: count}, nil
}

//查看我的收藏
func (s *FavoriteService) ListMyFavorites(userID uint) ([]model.Pet, error) {
	//1.从redis中读取收藏的宠物列表
	petIDs, err := s.favRepo.UserFavPetIDs(userID)
	if err != nil {
		return nil, err
	}
	if len(petIDs) == 0 {
		return []model.Pet{}, nil
	}

	//2.批量查看宠物信息
	var pets []model.Pet
	for _, id := range petIDs {
		pet, err := s.petRepo.FindByID(id)
		if err != nil {
			continue
		}
		pets = append(pets, *pet)
	}
	return pets, nil
}
