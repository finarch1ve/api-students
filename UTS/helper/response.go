package helper

import "github.com/gofiber/fiber/v2"

type Meta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	LastPage    int   `json:"last_page"`
}

func NewMeta(page, perPage int, total int64) Meta {
	last := int((total + int64(perPage) - 1) / int64(perPage))
	if last < 1 {
		last = 1
	}
	return Meta{CurrentPage: page, PerPage: perPage, Total: total, LastPage: last}
}

func Success(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func SuccessWithMeta(c *fiber.Ctx, message string, data interface{}, meta Meta) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}

func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}