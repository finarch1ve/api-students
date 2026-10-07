package main

import (
	"log"

	"siakad-mini/config"
	"siakad-mini/database"
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
}