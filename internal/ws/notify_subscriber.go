package ws

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/hacker4257/pet_charity/pkg/logger"
)

// StartNotifySubscriber 订阅 Redis 通知频道，推送给在线用户
func StartNotifySubscriber(hub *Hub) {
	go func() {
		ctx := context.Background()
		pubsub := database.RDB.PSubscribe(ctx, "notify:push:*")
		defer pubsub.Close()
		logger.Info("[notify] subscriber started")

		ch := pubsub.Channel()
		for msg := range ch {
			//格式notify:push:123
			parts := strings.Split(msg.Channel, ":")
			if len(parts) != 3 {
				continue
			}
			userID, err := strconv.ParseUint(parts[2], 10, 64)
			if err != nil {
				continue
			}
			//区分消息类型
			envelope, _ := json.Marshal(map[string]interface{}{
				"type": "notification",
				"data": json.RawMessage(msg.Payload),
			})
			hub.SendRawToUser(uint(userID), envelope)
		}
	}()
}
