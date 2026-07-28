package websocket

import (
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

type Hub struct {
	connections sync.Map
}

func NewHub() *Hub {
	return &Hub{}
}

func (h *Hub) HandleConnection() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		merchantID := c.Query("merchant_id")
		if merchantID == "" {
			c.WriteJSON(fiber.Map{"error": "missing merchant_id query parameter"})
			c.Close()
			return
		}

		h.connections.Store(merchantID, c)
		log.Printf("[ws] merchant %s connected", merchantID)

		c.WriteJSON(fiber.Map{
			"type":        "connected",
			"merchant_id": merchantID,
		})

		defer func() {
			h.connections.Delete(merchantID)
			log.Printf("[ws] merchant %s disconnected", merchantID)
			c.Close()
		}()

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}
		}
	})
}

func (h *Hub) PushToMerchant(merchantID string, payload interface{}) error {
	if conn, ok := h.connections.Load(merchantID); ok {
		return conn.(*websocket.Conn).WriteJSON(payload)
	}

	log.Printf("[ws] merchant %s not connected, skipping push", merchantID)
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
			log.Println("[ws] merchant cleanup routine running (placeholder)")
		}
	}()
}
