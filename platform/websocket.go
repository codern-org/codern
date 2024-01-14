package platform

import (
	"fmt"
	"sync"
	"time"

	"github.com/codern-org/codern/internal/constant"
	"github.com/gofiber/contrib/websocket"
)

type wsConnInfo struct {
	conn        *websocket.Conn
	createdTime time.Time
}

type WebSocketPayload struct {
	Channel string      `json:"channel"`
	Message interface{} `json:"message"`
}

type WebSocketChannelHandler func(message interface{})

type WebSocketHub struct {
	prometheus *Prometheus
	mu         sync.Mutex
	connPool   map[string][]wsConnInfo
	handlers   map[string]WebSocketChannelHandler
}

func NewWebSocketHub(prometheus *Prometheus) *WebSocketHub {
	return &WebSocketHub{
		prometheus: prometheus,
		connPool:   make(map[string][]wsConnInfo),
	}
}

// TODO: implement keep-alive
func (h *WebSocketHub) RegisterUser(userId string, conn *websocket.Conn) {
	// Call from fiber which need to be thread-safe
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.connPool[userId] == nil {
		h.connPool[userId] = make([]wsConnInfo, constant.MaxWebSocketConnPerUser)
		for i := range h.connPool[userId] {
			h.connPool[userId][i].createdTime = time.Now()
		}
		h.prometheus.GetUniqueActiveUserGauge().Inc()
	}

	// Find the oldest connection of the same user id
	oldest := &h.connPool[userId][0]
	for i := range h.connPool[userId] {
		if h.connPool[userId][i].conn == nil {
			oldest = &h.connPool[userId][i]
			break
		} else {
			if h.connPool[userId][i].createdTime.Before(oldest.createdTime) {
				oldest = &h.connPool[userId][i]
			}
		}
	}

	oldest.conn = conn
	oldest.createdTime = time.Now()
	h.prometheus.GetActiveUserGauge().Inc()
}

func (h *WebSocketHub) UnregisterUser(userId string, conn *websocket.Conn) {
	// Call from fiber which need to be thread-safe
	h.mu.Lock()
	defer h.mu.Unlock()

	for i := range h.connPool[userId] {
		if h.connPool[userId][i].conn == conn {
			h.connPool[userId][i].conn = nil
			break
		}
	}

	needToCleanUp := true
	for i := range h.connPool[userId] {
		if h.connPool[userId][i].conn != nil {
			needToCleanUp = false
			break
		}
	}
	if needToCleanUp {
		h.connPool[userId] = nil
		h.prometheus.GetUniqueActiveUserGauge().Dec()
	}

	h.prometheus.GetActiveUserGauge().Dec()
}

func (h *WebSocketHub) RegisterHandler(channel string, handler WebSocketChannelHandler) {
	h.handlers[channel] = handler
}

func (h *WebSocketHub) GetHandler(channel string) WebSocketChannelHandler {
	return h.handlers[channel]
}

func (h *WebSocketHub) SendMessage(userId string, channel string, message interface{}) error {
	connInfos := h.connPool[userId]
	if connInfos == nil {
		return fmt.Errorf("cannot get connection from user id %s", userId)
	}

	for i := range connInfos {
		if connInfos[i].conn == nil {
			continue
		}

		err := connInfos[i].conn.WriteJSON(WebSocketPayload{
			Channel: channel,
			Message: message,
		})
		if err != nil {
			return fmt.Errorf("cannot write json to websocket: %w", err)
		}
	}

	return nil
}

func (h *WebSocketHub) GetActiveCount() int {
	count := 0
	for i := range h.connPool {
		if len(h.connPool[i]) > 0 {
			for j := range h.connPool[i] {
				if h.connPool[i][j].conn != nil {
					count += 1
					continue
				}
			}
		}
	}
	return count
}
