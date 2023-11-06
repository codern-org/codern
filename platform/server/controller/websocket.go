package controller

import (
	"encoding/json"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/constant"
	"github.com/codern-org/codern/platform"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WebSocketController struct {
	wsHub *platform.WebSocketHub
}

func NewWebSocketController(
	wsHub *platform.WebSocketHub,
) *WebSocketController {
	return &WebSocketController{
		wsHub: wsHub,
	}
}

// Request to upgrade to the WebSocket protocol
func (c *WebSocketController) Upgrade(ctx *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		return ctx.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (c *WebSocketController) Portal() fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		user := conn.Locals(constant.UserCtxLocal).(*domain.User)
		c.wsHub.RegisterUser(user.Id, conn)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				// Error `websocket: close 1001 (going away)` = client disconnected
				c.wsHub.UnregisterUser(user.Id, conn)
				return
			}

			var payload platform.WebSocketPayload
			if err = json.Unmarshal(msg, &payload); err != nil {
				return
			}

			if handler := c.wsHub.GetHandler(payload.Channel); handler != nil {
				handler(payload.Message)
			}
		}
	})
}
