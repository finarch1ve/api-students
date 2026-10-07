package route

import (
	"siakad-mini/handler"
	"siakad-mini/helper"
	"siakad-mini/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api/v1")

	api.Get("/ping", func(c *fiber.Ctx) error {
		return helper.Success(c, fiber.StatusOK, "API SIAKAD Mini berjalan", nil)
	})

	auth := &handler.AuthHandler{DB: db}
	api.Post("/auth/login", middleware.LoginLimiter(), auth.Login)
	api.Get("/auth/me", middleware.AuthRequired(), auth.Me)
}