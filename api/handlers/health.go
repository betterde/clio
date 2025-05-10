package handlers

import (
	"github.com/betterde/clio/internal/response"
	"github.com/gofiber/fiber/v2"
)

func HealthCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(response.Success("Success", nil))
}
