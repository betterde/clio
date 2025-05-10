package handlers

import (
	"github.com/betterde/clio/internal/response"
	"github.com/gofiber/fiber/v2"
)

type SignInReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUp(ctx *fiber.Ctx) error {
	return ctx.JSON(response.Success("Success", nil))
}

func SignIn(ctx *fiber.Ctx) error {
	req := SignInReq{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.JSON(response.ValidationError("Failed to parse the request body.", err))
	}
	return ctx.JSON(response.Success("Success", nil))
}

func Callback(ctx *fiber.Ctx) error {
	return ctx.JSON(response.Success("Success", nil))
}
