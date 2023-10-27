package platform

import (
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketPayload struct {
	Channel string      `json:"channel"`
	Message interface{} `json:"message"`
}

type WebSocketChannelHandler func(message interface{})

type WebSocketHub struct {
	mu       sync.Mutex
	connPool map[string]*websocket.Conn
	handlers map[string]WebSocketChannelHandler
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		connPool: make(map[string]*websocket.Conn),
	}
}

func (h *WebSocketHub) RegisterUser(userId string, conn *websocket.Conn) {
	// Call from fiber which need to be thread-safe
	h.mu.Lock()
	h.connPool[userId] = conn
	h.mu.Unlock()
}

func (h *WebSocketHub) RegisterHandler(channel string, handler WebSocketChannelHandler) {
	h.handlers[channel] = handler
}

func (h *WebSocketHub) GetHandler(channel string) WebSocketChannelHandler {
	return h.handlers[channel]
}

func (h *WebSocketHub) SendMessage(userId string, channel string, message interface{}) error {
	conn := h.connPool[userId]
	if conn == nil {
		return fmt.Errorf("cannot get connection from user id %s", userId)
	}

	err := conn.WriteJSON(WebSocketPayload{
		Channel: channel,
		Message: message,
	})
	if err != nil {
		return fmt.Errorf("cannot write json to websocket: %w", err)
	}

	return nil
}
