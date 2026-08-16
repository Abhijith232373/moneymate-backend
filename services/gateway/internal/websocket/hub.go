package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/moneymate-2026/moneymate-backend/gateway/internal/proxy"
	"github.com/redis/go-redis/v9"
)

type Hub struct {
	rdb         *redis.Client
	connections sync.Map
}

func NewHub(rdb *redis.Client) *Hub {
	return &Hub{rdb: rdb}
}

func (h *Hub) HandleConnection(authClient proxy.AuthClient) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		token := c.Query("token")
		if token == "" {
			c.WriteJSON(fiber.Map{"error": "missing token query parameter"})
			c.Close()
			return
		}

		claims, err := authClient.VerifyAccessToken(context.Background(), token)
		if err != nil {
			c.WriteJSON(fiber.Map{"error": "invalid token"})
			c.Close()
			return
		}

		userID := claims.UserID

		sessionKey := fmt.Sprintf("session:%s", userID)
		h.rdb.Set(context.Background(), sessionKey, c.Conn.LocalAddr().String(), 0)

		h.connections.Store(userID, c)
		log.Printf("[ws] user %s connected", userID)

		c.WriteJSON(fiber.Map{
			"type":    "connected",
			"user_id": userID,
		})

		defer func() {
			h.connections.Delete(userID)
			h.rdb.Del(context.Background(), sessionKey)
			log.Printf("[ws] user %s disconnected", userID)
			c.Close()
		}()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}

			var incoming map[string]interface{}
			if err := json.Unmarshal(msg, &incoming); err == nil {
				payload := map[string]interface{}{
					"sender_id":     userID,
					"sender_type":   claims.Role,
					"receiver_id":   incoming["receiver_id"],
					"receiver_type": incoming["receiver_type"],
					"message":       incoming["message"],
				}
				payloadBytes, _ := json.Marshal(payload)
				h.rdb.Publish(context.Background(), "incoming_ws_messages", payloadBytes)
			}
		}
	})
}

func (h *Hub) PushToUser(userID string, payload interface{}) error {
	if conn, ok := h.connections.Load(userID); ok {
		return conn.(*websocket.Conn).WriteJSON(payload)
	}

	log.Printf("[ws] user %s not connected to this replica, skipping push", userID)
	return nil
}

func (h *Hub) Broadcast(payload interface{}) {
	h.connections.Range(func(key, value interface{}) bool {
		conn := value.(*websocket.Conn)
		conn.WriteJSON(payload)
		return true
	})
}

func (h *Hub) StartCleanupRoutine(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			log.Println("[ws] cleanup routine running (placeholder)")
		}
	}()
}

func (h *Hub) StartRedisSubscriber(ctx context.Context) {
	pubsub := h.rdb.Subscribe(ctx, "chat_messages")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		var chatMsg struct {
			SenderID   string `json:"sender_id"`
			ReceiverID string `json:"receiver_id"`
		}
		if err := json.Unmarshal([]byte(msg.Payload), &chatMsg); err == nil {
			var payload map[string]interface{}
			json.Unmarshal([]byte(msg.Payload), &payload)
			payload["type"] = "new_chat_message"
			
			h.PushToUser(chatMsg.ReceiverID, payload)
			h.PushToUser(chatMsg.SenderID, payload)
		}
	}
}
