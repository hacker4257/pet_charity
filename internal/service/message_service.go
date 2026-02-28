package service

import (
	"time"

	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/hacker4257/pet_charity/internal/repository"
	"github.com/hacker4257/pet_charity/internal/ws"
)

type MessageService struct {
	msgRepo  repository.MessageRepository
	userRepo repository.UserRepository
}

func NewMessageService(msgRepo repository.MessageRepository, userRepo repository.UserRepository) *MessageService {
	return &MessageService{msgRepo: msgRepo, userRepo: userRepo}
}

// SendPrivate 发送私聊消息
func (s *MessageService) SendPrivate(fromID, toID uint, content, msgType string) (*model.Message, error) {
	msg := &model.Message{
		FromUserID: fromID,
		ToUserID:   toID,
		Content:    content,
		MsgType:    msgType,
	}
	if err := s.msgRepo.Create(msg); err != nil {
		return nil, err
	}

	// 查发送者信息
	fromUser, _ := s.userRepo.FindByID(fromID)
	nickname := ""
	if fromUser != nil {
		nickname = fromUser.Nickname
		if nickname == "" {
			nickname = fromUser.Username
		}
	}

	// WebSocket 推送
	wsMsg := &ws.WsMessage{
		Type:       "chat",
		FromUserID: fromID,
		ToUserID:   toID,
		Content:    content,
		MsgType:    msgType,
		Nickname:   nickname,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}
	ws.GlobalHub.SendToUser(toID, wsMsg)
	ws.GlobalHub.SendToUser(fromID, wsMsg) // 自己也收一份

	return msg, nil
}

// SendToRoom 发送房间消息
func (s *MessageService) SendToRoom(fromID uint, roomID, content, msgType string) (*model.Message, error) {
	msg := &model.Message{
		FromUserID: fromID,
		RoomID:     roomID,
		Content:    content,
		MsgType:    msgType,
	}
	if err := s.msgRepo.Create(msg); err != nil {
		return nil, err
	}

	fromUser, _ := s.userRepo.FindByID(fromID)
	nickname := ""
	if fromUser != nil {
		nickname = fromUser.Nickname
		if nickname == "" {
			nickname = fromUser.Username
		}
	}

	wsMsg := &ws.WsMessage{
		Type:       "chat",
		FromUserID: fromID,
		RoomID:     roomID,
		Content:    content,
		MsgType:    msgType,
		Nickname:   nickname,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}
	ws.GlobalHub.BroadcastToRoom(roomID, wsMsg)

	return msg, nil
}

func (s *MessageService) ListPrivate(userA, userB uint, page, pageSize int) ([]model.Message, int64, error) {
	return s.msgRepo.ListPrivate(userA, userB, page, pageSize)
}

func (s *MessageService) ListByRoom(roomID string, page, pageSize int) ([]model.Message, int64, error) {
	return s.msgRepo.ListByRoom(roomID, page, pageSize)
}

func (s *MessageService) ListConversations(userID uint) ([]model.Message, error) {
	return s.msgRepo.ListConversations(userID)
}
