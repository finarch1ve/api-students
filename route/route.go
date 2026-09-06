package route

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"api-students/helper"
	"api-students/middleware"
)

// PingDB adalah kontrak minimal untuk mengecek koneksi database.
type PingDB interface {
	Ping(ctx context.Context) error
}

// StudentHandlers menampung fungsi-fungsi handler mahasiswa.
type StudentHandlers struct {
	List    fiber.Handler
	Get     fiber.Handler
	Create  fiber.Handler
	Replace fiber.Handler
	Patch   fiber.Handler
	Delete  fiber.Handler
}

// Register mendaftarkan seluruh route API ke aplikasi Fiber.
func Register(app *fiber.App, pool PingDB, h StudentHandlers) {
	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			return helper.Fail(c, fiber.StatusServiceUnavailable, "database tidak dapat dihubungi")
		}
		return helper.Success(c, fiber.StatusOK, "server dan database berjalan", nil)
	})

	s := api.Group("/students", middleware.RequireJSON)
	s.Get("/", h.List)
	s.Get("/:id", h.Get)
	s.Post("/", h.Create)
	s.Put("/:id", h.Replace)
	s.Patch("/:id", h.Patch)
	s.Delete("/:id", h.Delete)

	app.Use(func(c *fiber.Ctx) error {
		return helper.Fail(c, fiber.StatusNotFound, "endpoint tidak ditemukan")
	})
}