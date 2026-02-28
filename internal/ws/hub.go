package ws

import (
	"encoding/json"
	"sync"
)

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

//发送私聊信息
func (h *Hub) SendToUser(userID uint, msg *WsMessage) {
	data, _ := json.Marshal(msg)
	h.mu.RLock()
	defer h.mu.RUnlock()
	if client, ok := h.clients[userID]; ok {
		select {
		case client.Send <- data:
		default:
			//缓存满了，丢弃
		}
	}
}

//广播房间
func (h *Hub) BroadcastToRoom(roomID string, msg *WsMessage) {
	data, _ := json.Marshal(msg)
	h.mu.RLock()
	defer h.mu.RUnlock()
	if members, ok := h.rooms[roomID]; ok {
		for _, client := range members {
			select {
			case client.Send <- data:
			default:
			}
		}
	}
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

// 发送消息
func (h *Hub) SendRawToUser(userID uint, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if client, ok := h.clients[userID]; ok {
		select {
		case client.Send <- data:
		default:
			// 缓冲满，丢弃
		}
	}
}
