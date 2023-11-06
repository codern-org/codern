package server

import (
	"encoding/json"
	"errors"

	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/internal/constant"
	"github.com/codern-org/codern/platform/server/response"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func errorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		resStatus := fiber.StatusInternalServerError
		requestId := ctx.Locals(constant.RequestIdCtxLocal).(string)

		var domainError *errs.DomainError
		if errors.As(err, &domainError) {
			status, ok := response.DomainErrCodeToHttpStatus[domainError.Code]
			if ok {
				resStatus = status
			}
		}

		// Handle websocket
		if websocket.IsWebSocketUpgrade(ctx) {
			msg, err := json.Marshal(fiber.Map{"error": err})
			if err != nil {
				logger.Error("Cannot handle websocket connection", zap.NamedError("error", err))
				return err
			}

			if err := wsErrHandler(string(msg))(ctx); err != nil {
				logger.Error("Cannot handle websocket connection", zap.NamedError("error", err))
				return err
			}
			return nil
		}

		// Mostly JSON format
		if jsonErr := response.NewErrorResponse(ctx, resStatus, err); jsonErr != nil {
			logger.Error(
				"Cannot marshal json response",
				zap.String("request_id", requestId),
				zap.NamedError("json_error", jsonErr),
				zap.NamedError("error", err),
			)
			// Fallback to string format if cannot marshal JSON
			return ctx.Status(500).SendString(err.Error())
		}

		return nil
	}
}

var wsErrHandler = func(message string) fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		closeMsg := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, message)
		conn.WriteMessage(websocket.CloseMessage, closeMsg)
	})
}
