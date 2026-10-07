package middleware

import (
	"strings"

	"siakad-mini/helper"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		parts := strings.SplitN(c.Get("Authorization"), " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			return helper.Unauthorized("Token tidak ditemukan atau format salah")
		}
		claims, err := helper.ParseToken(strings.TrimSpace(parts[1]))
		if err != nil {
			return helper.Unauthorized("Token tidak valid atau kedaluwarsa")
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}
		return helper.Forbidden("Anda tidak memiliki akses ke resource ini")
	}
}