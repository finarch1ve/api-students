package main

import (
	"log"
	"os"

	"siakad-mini/config"
	"siakad-mini/database"
	"siakad-mini/helper"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Gagal koneksi database: ", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal("Gagal migrasi: ", err)
	}
	log.Println("Migrasi selesai")

	if err := database.Seed(db); err != nil {
		log.Fatal("Gagal seeding: ", err)
	}
	log.Println("Seeder selesai")

	app := fiber.New(fiber.Config{ErrorHandler: config.ErrorHandler})

	app.Get("/api/v1/ping", func(c *fiber.Ctx) error {
		return helper.Success(c, fiber.StatusOK, "API SIAKAD Mini berjalan", nil)
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}