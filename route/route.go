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

// StudentHandlers menampung fungsi-fungsi handler mahasiswa dan prestasi.
type StudentHandlers struct {
	List             fiber.Handler
	Get              fiber.Handler
	Create           fiber.Handler
	Replace          fiber.Handler
	Patch            fiber.Handler
	Delete           fiber.Handler
	GetPrestasiByNIM fiber.Handler
}

// AuthHandlers menampung fungsi-fungsi handler autentikasi.
type AuthHandlers struct {
	Register fiber.Handler
	Login    fiber.Handler
	Refresh  fiber.Handler
	Logout   fiber.Handler
	Me       fiber.Handler
}

// Dependencies mengumpulkan semua yang dibutuhkan Register dalam satu struct.
type Dependencies struct {
	Pool     PingDB
	JWT      *helper.JWTManager
	Students StudentHandlers
	Auth     AuthHandlers
}

// Register mendaftarkan seluruh route API ke aplikasi Fiber.
func Register(app *fiber.App, deps Dependencies) {
	api := app.Group("/api/v1")

	// --- publik ---
	api.Get("/health", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()
		if err := deps.Pool.Ping(ctx); err != nil {
			return helper.Fail(c, fiber.StatusServiceUnavailable, "database tidak dapat dihubungi")
		}
		return helper.Success(c, fiber.StatusOK, "server dan database berjalan", nil)
	})

	// --- autentikasi ---
	auth := api.Group("/auth", middleware.RequireJSON)
	auth.Post("/register", deps.Auth.Register)
	auth.Post("/login", middleware.LoginRateLimiter(), deps.Auth.Login)
	auth.Post("/refresh", deps.Auth.Refresh)
	auth.Post("/logout", deps.Auth.Logout)
	auth.Get("/me", middleware.RequireAuth(deps.JWT), deps.Auth.Me)

	// --- wajib membawa access token ---
	s := api.Group("/students", middleware.RequireJSON, middleware.RequireAuth(deps.JWT))
	s.Get("/", deps.Students.List)
	s.Get("/:id", deps.Students.Get)
	s.Post("/", deps.Students.Create)
	s.Put("/:id", deps.Students.Replace)
	s.Patch("/:id", deps.Students.Patch)
	s.Delete("/:id", deps.Students.Delete)
	s.Get("/nim/:nim/prestasi", deps.Students.GetPrestasiByNIM)

	app.Use(func(c *fiber.Ctx) error {
		return helper.Fail(c, fiber.StatusNotFound, "endpoint tidak ditemukan")
	})
}