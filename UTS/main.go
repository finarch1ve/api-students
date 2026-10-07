package main

import (
	"log"
	"os"

	"siakad-mini/config"
	"siakad-mini/database"
	"siakad-mini/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET belum diisi di .env")
	}

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
	route.Setup(app, db)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}