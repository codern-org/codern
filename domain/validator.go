package domain

import "github.com/gofiber/fiber/v2"

type PayloadValidator interface {
	Validate(payload interface{}, ctx *fiber.Ctx) (bool, error)
}
