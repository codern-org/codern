package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewLogger(logger *zap.Logger, influxdb domain.InfluxDb) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		startTime := time.Now()
		err := ctx.Next()
		executionTime := time.Since(startTime)

		method := ctx.Method()
		path := ctx.Path()
		statusCode := ctx.Response().StatusCode()
		ip := ctx.IP()
		requestId := ctx.Locals("requestid").(string)
		userAgent := ctx.Context().UserAgent()

		influxDbErr := influxdb.WritePoint(
			"httpRequest",
			map[string]string{
				"method":     method,
				"path":       path,
				"statusCode": strconv.Itoa(statusCode),
				"ipAddress":  ip,
			},
			map[string]interface{}{
				"executionTime": executionTime.Nanoseconds(),
			},
		)
		if influxDbErr != nil {
			logger.Error("HTTP Request Measurement", zap.Error(influxDbErr))
			return ctx.Status(fiber.StatusInternalServerError).JSON(response.GenericErrorResponse{
				Code:    domain.ErrLoggingError,
				Message: "Internal logging error",
			})
		}

		logger.Info(
			fmt.Sprintf("Request %s %s %d", method, path, statusCode),
			zap.String("request_id", requestId),
			zap.String("ip_address", ip),
			zap.String("user_agent", string(userAgent)),
			zap.String("execution_time", executionTime.String()),
		)

		return err
	}
}
