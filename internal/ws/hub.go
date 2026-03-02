package ws

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"sync"

	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/hacker4257/pet_charity/pkg/logger"
)

// 实例唯一标识，用于 Redis Pub/Sub 去重
var instanceID string

func init() {
	b := make([]byte, 8)
	rand.Read(b)
	instanceID = hex.EncodeToString(b)
}

type Client struct {
	UserID uint
	Conn   interface {
		WriteMessage(int, []byte) error
		Close() error
	}
	Send chan []byte
}

// 消息体
type WsMessage struct {
	Type       string `json:"type"`
	FromUserID uint   `json:"from_user_id"`
	ToUserID   uint   `json:"to_user_id"`
	RoomID     string `json:"room_id,omitempty"`
	Content    string `json:"content"`
	MsgType    string `json:"msg_type"`
	Nickname   string `json:"nickname,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
}

// redisEnvelope Redis 传输的信封，带实例 ID 用于去重
type redisEnvelope struct {
	Instance string    `json:"ins"`
	Msg      WsMessage `json:"msg"`
}

// Redis 频道名
const (
	channelChat = "ws:chat"
	channelRoom = "ws:room"
)

// hub管理所有连接
type Hub struct {
	mu         sync.RWMutex
	clients    map[uint]*Client
	rooms      map[string]map[uint]*Client
	Register   chan *Client
	Unregister chan *Client
}

var GlobalHub *Hub

func NewHub() *Hub {
	h := &Hub{
		clients:    make(map[uint]*Client),
		rooms:      make(map[string]map[uint]*Client),
		Register:   make(chan *Client, 256),
		Unregister: make(chan *Client, 256),
	}
	GlobalHub = h
	return h
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				close(client.Send)
				delete(h.clients, client.UserID)
			}
			//移除房间
			for roomID, members := range h.rooms {
				delete(members, client.UserID)
				if len(members) == 0 {
					delete(h.rooms, roomID)
				}
			}
			h.mu.Unlock()
		}
	}
}

// sendLocal 本地投递，返回是否成功
func (h *Hub) sendLocal(userID uint, data []byte) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if client, ok := h.clients[userID]; ok {
		select {
		case client.Send <- data:
		default:
		}
		return true
	}
	return false
}

// publishToRedis 发送带实例标识的信封到 Redis
func publishToRedis(channel string, msg *WsMessage) {
	env := redisEnvelope{Instance: instanceID, Msg: *msg}
	data, _ := json.Marshal(env)
	database.RDB.Publish(context.Background(), channel, data)
}

// SendToUser 发送私聊：本地优先，找不到则通过 Redis 转发
func (h *Hub) SendToUser(userID uint, msg *WsMessage) {
	data, _ := json.Marshal(msg)
	if h.sendLocal(userID, data) {
		return
	}
	// 本地没有，发到 Redis 让其他实例投递
	publishToRedis(channelChat, msg)
}

// BroadcastToRoom 广播房间：本地投递 + Redis 转发给其他实例
func (h *Hub) BroadcastToRoom(roomID string, msg *WsMessage) {
	data, _ := json.Marshal(msg)
	// 本地投递
	h.mu.RLock()
	if members, ok := h.rooms[roomID]; ok {
		for _, client := range members {
			select {
			case client.Send <- data:
			default:
			}
		}
	}
	h.mu.RUnlock()
	// 发到 Redis，让其他实例投递给各自的房间成员
	publishToRedis(channelRoom, msg)
}

//加入房间
func (h *Hub) JoinRoom(roomID string, userID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[uint]*Client)
	}
	if client, ok := h.clients[userID]; ok {
		h.rooms[roomID][userID] = client
	}
}

//离开房间
func (h *Hub) LeaveRoom(roomID string, userID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if members, ok := h.rooms[roomID]; ok {
		delete(members, userID)
	}
}

// SendRawToUser 发送原始字节（通知等场景）
func (h *Hub) SendRawToUser(userID uint, data []byte) {
	if h.sendLocal(userID, data) {
		return
	}
	// 包装成 WsMessage 再转发
	var msg WsMessage
	if json.Unmarshal(data, &msg) == nil {
		publishToRedis(channelChat, &msg)
	}
}

// StartChatSubscriber 订阅 Redis 聊天频道，接收其他实例转发的消息
func StartChatSubscriber(hub *Hub) {
	go func() {
		ctx := context.Background()
		sub := database.RDB.Subscribe(ctx, channelChat, channelRoom)
		defer sub.Close()
		logger.Info("[ws] chat subscriber started")

		for msg := range sub.Channel() {
			var env redisEnvelope
			if json.Unmarshal([]byte(msg.Payload), &env) != nil {
				continue
			}
			// 跳过自己发出的消息，避免重复投递
			if env.Instance == instanceID {
				continue
			}

			wsData, _ := json.Marshal(env.Msg)

			switch msg.Channel {
			case channelChat:
				// 私聊：尝试本地投递
				if env.Msg.ToUserID > 0 {
					hub.sendLocal(env.Msg.ToUserID, wsData)
				}
			case channelRoom:
				// 房间消息：投递给本地房间成员
				hub.mu.RLock()
				if members, ok := hub.rooms[env.Msg.RoomID]; ok {
					for _, client := range members {
						select {
						case client.Send <- wsData:
						default:
						}
					}
				}
				hub.mu.RUnlock()
			}
		}
	}()
}
