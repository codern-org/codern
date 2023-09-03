package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/response"
	"github.com/codern-org/codern/platform"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewLogger(logger *zap.Logger, influxdb *platform.InfluxDb) fiber.Handler {
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
			},
			map[string]interface{}{
				"executionTime": executionTime.Nanoseconds(),
			},
		)
		if influxDbErr != nil {
			logger.Error("HTTP Request Measurement", zap.Error(influxDbErr))
			return response.NewErrorResponse(
				ctx,
				fiber.StatusInternalServerError,
				domain.NewError(domain.ErrLoggingError, "internal logging error"),
			)
		}

		logger.Info(
			fmt.Sprintf("Request %s %s %d", method, path, statusCode),
			zap.String("request_id", requestId),
			zap.String("ip_address", ip),
			zap.String("user_agent", string(userAgent)),
			zap.String("execution_time", executionTime.String()),
		)

		if err != nil {
			logger.Error("Internal server error", zap.Error(err))
		}

		internalError := ctx.Locals("internal_error")
		if internalError != nil {
			logger.Error("Internal server error", zap.Error(internalError.(error)))
		}

		return err
	}
}
