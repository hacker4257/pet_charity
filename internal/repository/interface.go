package repository

import (
	"time"

	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/redis/go-redis/v9"
)

type UserRepository interface {
	FindByID(id uint) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByPhone(phone string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	UpdateFields(id uint, fields map[string]interface{}) error
	Count() (int64, error)
	List(page, pageSize int, role string) ([]model.User, int64, error)
	UpdateRole(id uint, role string) error
	UpdateStatus(id uint, status string) error
}

type OrgRepository interface {
	FindByID(id uint) (*model.Organization, error)
	FindByUserID(userID uint) (*model.Organization, error)
	Create(org *model.Organization) error
	UpdateFields(id uint, fields map[string]interface{}) error
	ListApproved(page, pageSize int) ([]model.Organization, int64, error)
	ListPending(page, pageSize int) ([]model.Organization, int64, error)
	FindNearby(lng, lat float64, radiusKm float64, limit int) ([]model.Organization, error)
	ApproveWithTX(orgID uint, userID uint) error
	CountByStatus(status string) (int64, error)
}

type PetRepository interface {
	Create(pet *model.Pet) error
	FindByID(id uint) (*model.Pet, error)
	UpdateFields(id uint, fields map[string]interface{}) error
	Delete(id uint) error
	List(filter PetFilter, page, pageSize int) ([]model.Pet, int64, error)
	CreateImage(image *model.PetImage) error
	DeleteImage(imageID uint) error
	FindImageByID(imageID uint) (*model.PetImage, error)
	PublicStats() map[string]int64
}

type AdoptionRepository interface {
	Create(adoption *model.Adoption) error
	FindByID(id uint) (*model.Adoption, error)
	UpdateFields(id uint, fields map[string]interface{}) error
	FindPendingByUserAndPet(userID, petID uint) (*model.Adoption, error)
	FindApprovedByUserAndPet(userID, petID uint) (*model.Adoption, error)
	ListByUser(userID uint, page, pageSize int) ([]model.Adoption, int64, error)
	ListByOrg(orgID uint, status string, page, pageSize int) ([]model.Adoption, int64, error)
}

type DonationRepository interface {
	Create(donation *model.Donation) error
	FindByID(id uint) (*model.Donation, error)
	FindByTradeNo(tradeNo string) (*model.Donation, error)
	UpdateFields(id uint, fields map[string]interface{}) error
	ListByUser(userID uint, page, pageSize int) ([]model.Donation, int64, error)
	ListPublic(targetType string, targetID uint, page, pageSize int) ([]model.Donation, int64, error)
	Stats() (*DonationStats, error)
	UpdateFieldsWithStatus(id uint, expectStatus string, fields map[string]interface{}) (int64, error)
	FindPendingByUser(userID uint, targetType string, targetID uint) (*model.Donation, error)
	AddExpireTask(tradeNo string, expireAt time.Time) error
	GetExpiredTasks() ([]string, error)
	RemoveExpireTask(tradeNo string) error
}

type RescueRepository interface {
	Create(rescue *model.Rescue) error
	FindByID(id uint) (*model.Rescue, error)
	UpdateFields(id uint, fields map[string]interface{}) error
	List(filter RescueFilter, page, pageSize int) ([]model.Rescue, int64, error)
	ListForMap() ([]model.Rescue, error)
	FindNearby(lng, lat float64, radiusKm float64, limit int) ([]model.Rescue, error)
	CreateImage(image *model.RescueImage) error
	CreateFollow(follow *model.RescueFollow) error
	CreateClaim(claim *model.RescueClaim) error
	FindClaimByRescueAndOrg(rescueID, orgID uint) (*model.RescueClaim, error)
	FindActiveClaimByRescue(rescueID uint) (*model.RescueClaim, error)
	UpdateClaimFields(id uint, fields map[string]interface{}) error
}

type SmsRepository interface {
	SaveCode(purpose, phone, code string) error
	GetCode(purpose, phone string) (string, error)
	DeleteCode(purpose, phone string) error
	CheckThrottle(phone string) (bool, error)
	SetThrottle(phone string) error
	IncrDaily(phone string) (int64, error)
}

type MessageRepository interface {
	Create(msg *model.Message) error
	ListPrivate(userA, userB uint, page, pageSize int) ([]model.Message, int64, error)
	ListByRoom(roomID string, page, pageSize int) ([]model.Message, int64, error)
	ListConversations(userID uint) ([]model.Message, error)
}

type ActivityRepository interface {
	CreateLog(log *model.UserActivityLogs) error
	AddScore(userID uint, points int) error
	RedisAddScore(userID uint, points int) error
	GetTopN(page, pageSize int) ([]redis.Z, error)
	GetUserRank(userID uint) (rank int64, score float64, err error)
	GetUsersByIDs(ids []uint) ([]model.User, error)
}

type FavoriteRepository interface {
	Add(userID, petID uint) error
	Remove(userID, petID uint) error
	IsFavorited(userID, petID uint) (bool, error)
	PetFavCount(petID uint) (int64, error)
	UserFavPetIDs(userID uint) ([]uint, error)
}

type CacheRepository interface {
	SetUser(userID uint, nickname, avatar string) error
	GetUser(userID uint) (*CachedUser, error)
	DeleteUser(userID uint) error
}

// TransactionManager 事务管理器抽象
type TransactionManager interface {
	// Transaction 在事务中执行 fn，fn 返回 error 则回滚，否则提交
	Transaction(fn func(tx TransactionContext) error) error
}

// TransactionContext 事务上下文，能获取事务版本的 repo
type TransactionContext interface {
	AdoptionRepo() AdoptionRepository
	PetRepo() PetRepository
	DonationRepo() DonationRepository
}

type NotificationRepository interface {
	Create(notification *model.Notification) error
	ListByUser(userID uint, page, pageSize int) ([]model.Notification, int64, error)
	MarkRead(userID, id uint) error
	MarkAllRead(userID uint) error
	GetUnreadCount(userID uint) (int64, error)
	IncrUnread(userID uint) error
	DecrUnread(userID uint) error
	Publish(userID uint, data []byte) error
	ResetUnread(userID uint) error
}


type PetDiaryRepository interface {
	Create(diary *model.PetDiary) error
	FindByID(id uint) (*model.PetDiary, error)
	Update(id uint, fields map[string]interface{}) error
	Delete(id uint) error
	ListByPet(petID uint, page, pageSize int) ([]model.PetDiary, int64, error)
	ListByUser(userID uint, page, pageSize int) ([]model.PetDiary, int64, error)
	ListPublic(page, pageSize int) ([]model.PetDiary, int64, error)

	// 图片
	CreateImage(image *model.DiaryImage) error
	DeleteImage(imageID uint) error
	FindImageByID(imageID uint) (*model.DiaryImage, error)

	// 点赞
	ToggleLike(diaryID, userID uint) (bool, error)
	CountLikes(diaryID uint) (int64, error)
	IsLiked(diaryID, userID uint) (bool, error)
}
