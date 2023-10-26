package middleware

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

		// Manually call error handler
		if chainErr != nil {
			if err := ctx.App().ErrorHandler(ctx, chainErr); err != nil {
				ctx.SendStatus(fiber.StatusInternalServerError)
			}
		}
		statusCode := ctx.Response().StatusCode()

		pathType, ok := ctx.Locals("pathType").(string)
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
		}
		logMessage := fmt.Sprintf("Request %s %s %d", method, path, statusCode)

		// Log with info level if status code is 2xx (successfull) or 1xx (informational)
		if strings.HasPrefix(fmt.Sprint(statusCode), "2") || strings.HasPrefix(fmt.Sprint(statusCode), "1") {
			if path != "/health" { // Ignore health check path
				logger.Info(logMessage, logFields...)
			}
		} else {
			logFields = append(logFields, zap.Error(chainErr))
			logger.Error(logMessage, logFields...)
		}

		return nil
	}
}
