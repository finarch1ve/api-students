package middleware

import (
	"time"

	"siakad-mini/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func LoginLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:                    5,
		Expiration:             time.Minute,
		SkipSuccessfulRequests: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return helper.NewAppError(fiber.StatusTooManyRequests, "Terlalu banyak percobaan login, coba lagi nanti")
		},
	})
}