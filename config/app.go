package config

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"api-students/helper"
	"api-students/middleware"
	"api-students/route"
)

// NewApp merakit aplikasi Fiber: config dasar, middleware, lalu route.
func NewApp(logger *slog.Logger, deps route.Dependencies) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      GetEnv("APP_NAME", "API Students - Tugas Mandiri Pertemuan 5"),
		ErrorHandler: newErrorHandler(logger),
		BodyLimit:    1 * 1024 * 1024, // 1 MB
	})

	middleware.Register(app, logger)
	route.Register(app, deps)

	return app
}

func newErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		status := fiber.StatusInternalServerError
		message := "terjadi kesalahan pada server"
		if e, ok := err.(*fiber.Error); ok {
			status = e.Code
			message = e.Message
		}
		logger.Error("unhandled_error",
			slog.String("path", c.Path()),
			slog.Int("status", status),
			slog.String("error", err.Error()),
		)
		return helper.Fail(c, status, message)
	}
}