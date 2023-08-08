package domain

import "github.com/gofiber/fiber/v2"

type PayloadValidator interface {
	ValidateAuth(ctx *fiber.Ctx) (string, error)
	ValidateBody(payload interface{}, ctx *fiber.Ctx) (bool, error)
}
