package service

import (
	"encoding/json"

	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/hacker4257/pet_charity/internal/repository"
	"github.com/hacker4257/pet_charity/pkg/event"
)

type NotificationService struct {
	notifyRepo repository.NotificationRepository
}

func NewNotificationService(notifyRepo repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		notifyRepo: notifyRepo,
	}
}

// 通知列表
func (s *NotificationService) List(userID uint, page, pageSize int) ([]model.Notification,
	int64, error) {
	return s.notifyRepo.ListByUser(userID, page, pageSize)
}

// 未读数
func (s *NotificationService) UnreadCount(userID uint) (int64, error) {
	return s.notifyRepo.GetUnreadCount(userID)
}

// 标记单条已读
func (s *NotificationService) MarkRead(userID, notificationID uint) error {
	if err := s.notifyRepo.MarkRead(userID, notificationID); err != nil {
		return err
	}
	s.notifyRepo.DecrUnread(userID)
	return nil
}

// 全部已读
func (s *NotificationService) MarkAllRead(userID uint) error {
	if err := s.notifyRepo.MarkAllRead(userID); err != nil {
		return err
	}
	s.notifyRepo.ResetUnread(userID)
	return nil
}

//发送通知
func (s *NotificationService) Send(userID uint, typ, title, content string, relatedID uint) error {

	//1.写入MySQL
	notification := &model.Notification{
		UserID:    userID,
		Type:      typ,
		Title:     title,
		Content:   content,
		RelatedID: relatedID,
	}

	if err := s.notifyRepo.Create(notification); err != nil {
		return err
	}

	//2.未读加1
	s.notifyRepo.IncrUnread(userID)

	//3.推送给在线用户
	push, _ := json.Marshal(map[string]interface{}{
		"id":         notification.ID,
		"type":       typ,
		"title":      title,
		"content":    content,
		"related_id": relatedID,
	})
	s.notifyRepo.Publish(userID, push)

	return nil
}

//事件->通知的映射
func (s *NotificationService) RegisterHook() {
	event.Subscribe(func(e event.Event) {
		switch e.Action {
		case "adoption_approved":
			s.Send(e.UserID, "adoption_approved",
				"领养申请已通过",
				"您的领养申请已通过，请联系机构办理手续",
				e.RelatedID)
		case "adoption_rejected":
			s.Send(e.UserID, "adoption_rejected",
				"领养申请未通过",
				"很遗憾，您的领养申请未通过审核",
				e.RelatedID)

		case "donation":
			s.Send(e.UserID, "donation_success",
				"捐赠成功",
				"感谢您的爱心捐赠，已成功到账",
				e.RelatedID)

		case "rescue_claimed":
			s.Send(e.UserID, "rescue_claimed",
				"救助已被认领",
				"您发起的救助已被机构认领，正在处理中",
				e.RelatedID)

		}
	})
}
