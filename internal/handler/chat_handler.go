package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hacker4257/pet_charity/internal/middleware"
	"github.com/hacker4257/pet_charity/internal/service"
	"github.com/hacker4257/pet_charity/internal/ws"
	"github.com/hacker4257/pet_charity/pkg/event"
	"github.com/hacker4257/pet_charity/pkg/logger"
	"github.com/hacker4257/pet_charity/pkg/response"
	"github.com/hacker4257/pet_charity/pkg/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type ChatHandler struct {
	msgService *service.MessageService
	hub        *ws.Hub
}

func NewChatHandler(msgService *service.MessageService, hub *ws.Hub) *ChatHandler {
	return &ChatHandler{msgService: msgService, hub: hub}
}

//websocket 连接入口
func (h *ChatHandler) HandleWS(c *gin.Context) {
	userID := middleware.GetUserID(c)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Warnf("ws upgrade error: %v", err)
		return
	}

	client := &ws.Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte),
	}
	h.hub.Register <- client
	go h.writePump(client, conn)
	go h.readPump(client, conn, userID)
}

//读取信息
func (h *ChatHandler) readPump(client *ws.Client, conn *websocket.Conn, userID uint) {
	defer func() {
		h.hub.Unregister <- client
		conn.Close()
	}()

	conn.SetReadLimit(64 * 1024) // 64KB per message
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var msg ws.WsMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		switch msg.Type {
		case "chat":
			if msg.RoomID != "" {
				h.msgService.SendToRoom(userID, msg.RoomID, msg.Content, "text")
			} else if msg.ToUserID > 0 {
				h.msgService.SendPrivate(userID, msg.ToUserID, msg.Content, "text")
			}
			event.Publish("chat", userID, 0)
		case "join_room":
			h.hub.JoinRoom(msg.RoomID, userID)
		case "leave_room":
			h.hub.LeaveRoom(msg.RoomID, userID)
		case "ping":
			//保活，不做处理

		}
	}
}

//向客户端写消息
func (h *ChatHandler) writePump(client *ws.Client, conn *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case msg, ok := <-client.Send:
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

//rest 获取会话列表
func (h *ChatHandler) Conversations(c *gin.Context) {
	userID := middleware.GetUserID(c)
	msgs, err := h.msgService.ListConversations(userID)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}
	response.Success(c, msgs)
}

// REST: 获取私聊历史
func (h *ChatHandler) PrivateHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	targetID, _ := strconv.ParseUint(c.Param("userId"), 10, 64)
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	msgs, total, err := h.msgService.ListPrivate(userID, uint(targetID), page, pageSize)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}
	response.SuccessWithPage(c, msgs, total, page, pageSize)
}

// REST: 获取房间历史
func (h *ChatHandler) RoomHistory(c *gin.Context) {
	roomID := c.Param("roomId")
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	msgs, total, err := h.msgService.ListByRoom(roomID, page, pageSize)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}
	response.SuccessWithPage(c, msgs, total, page, pageSize)
}
