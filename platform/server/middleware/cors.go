package middleware

import (
	"github.com/codern-org/codern/internal/constant"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var Cors = cors.New(cors.Config{
	AllowCredentials: true,
	AllowOriginsFunc: func(origin string) bool {
		return constant.IsDevelopment
	},
})
