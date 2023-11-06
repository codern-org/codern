package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/codern-org/codern/internal/constant"
	"github.com/codern-org/codern/platform"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(logger *zap.Logger, influxdb *platform.InfluxDb) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		startTime := time.Now()
		chainErr := ctx.Next()
		executionTime := time.Since(startTime)

		method := ctx.Method()
		path := ctx.Path()
		ip := ctx.IP()
		requestId := ctx.Locals("requestid").(string)
		userAgent := ctx.Context().UserAgent()

		user := GetUserFromCtx(ctx)

		// Manually call error handler
		if chainErr != nil {
			if err := ctx.App().ErrorHandler(ctx, chainErr); err != nil {
				// Maybe deadcode
				ctx.SendStatus(fiber.StatusInternalServerError)
				logger.Error("Unable to handle error", zap.Error(err))
			}
		}
		statusCode := ctx.Response().StatusCode()

		pathType, ok := ctx.Locals(constant.PathTypeCtxLocal).(string)
		if !ok {
			pathType = "unknown"
		}

		influxdb.WritePoint(
			"httpRequest",
			map[string]string{
				"method":     method,
				"pathType":   pathType,
				"statusCode": strconv.Itoa(statusCode),
			},
			map[string]interface{}{
				"path":          path,
				"executionTime": executionTime.Nanoseconds(),
				"ipAddress":     ip,
				"userAgent":     string(userAgent),
			},
		)

		logFields := []zapcore.Field{
			zap.String("request_id", requestId),
			zap.String("ip_address", ip),
			zap.String("user_agent", string(userAgent)),
			zap.String("execution_time", executionTime.String()),
			zap.Error(chainErr),
		}
		logMessage := fmt.Sprintf("Request %s %s %d", method, path, statusCode)

		if user != nil {
			logFields = append(logFields, zap.String("user_id", user.Id))
		}

		if chainErr == nil {
			if path != "/health" { // Ignore health check path
				logger.Info(logMessage, logFields...)
			}
		} else {
			logger.Error(logMessage, logFields...)
		}

		return nil
	}
}
