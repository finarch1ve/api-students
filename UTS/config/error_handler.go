package config

import (
	"errors"
	"log"
	"os"

	"siakad-mini/helper"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *helper.AppError
	if errors.As(err, &appErr) {
		body := fiber.Map{"success": false, "message": appErr.Message}
		if appErr.Errors != nil {
			body["errors"] = appErr.Errors
		}
		return c.Status(appErr.Status).JSON(body)
	}

	var fe *fiber.Error
	if errors.As(err, &fe) {
		return c.Status(fe.Code).JSON(fiber.Map{"success": false, "message": fe.Message})
	}

	log.Printf("ERROR %s %s: %v", c.Method(), c.Path(), err)
	msg := "Terjadi kesalahan pada server"
	if os.Getenv("APP_ENV") != "production" {
		msg = err.Error()
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
		"message": msg,
	})
}