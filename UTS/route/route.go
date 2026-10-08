package route

import (
	"siakad-mini/handler"
	"siakad-mini/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api/v1")

	auth := &handler.AuthHandler{DB: db}
	api.Post("/auth/login", middleware.LoginLimiter(), auth.Login)
	api.Get("/auth/me", middleware.AuthRequired(), auth.Me)

	st := &handler.StudentHandler{DB: db}
	students := api.Group("/students", middleware.AuthRequired())
	students.Get("/", middleware.RequireRole("admin"), st.List)
	students.Post("/", middleware.RequireRole("admin"), st.Create)
	students.Get("/:id", st.Show)
	students.Put("/:id", middleware.RequireRole("admin"), st.Update)
	students.Delete("/:id", middleware.RequireRole("admin"), st.Delete)
}