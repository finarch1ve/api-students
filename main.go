package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-students/app/repository"
	"api-students/app/service"
	"api-students/config"
	"api-students/database"
	"api-students/helper"
	"api-students/route"
)

const minSecretLength = 32

func main() {
	// 1. Konfigurasi dan logger
	config.LoadEnv()
	logger := config.NewLogger()

	// JWT_SECRET diperiksa sebelum server menyala.
	jwtSecret := config.GetEnv("JWT_SECRET", "")
	if len(jwtSecret) < minSecretLength {
		logger.Error("JWT_SECRET tidak diisi atau terlalu pendek", slog.Int("minimal_karakter", minSecretLength))
		os.Exit(1)
	}

	// 2. Database
	pool, err := database.NewPool(context.Background())
	if err != nil {
		logger.Error("gagal terhubung ke database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	jwtManager := helper.NewJWTManager(
		jwtSecret,
		config.GetEnv("JWT_ISSUER", "api-students"),
		time.Duration(config.GetEnvInt("JWT_ACCESS_TTL_MINUTES", 15))*time.Minute,
	)

	// 3. Perakitan: repository -> service -> handler
	studentRepository := repository.NewStudentRepository(pool)
	studentService := service.NewStudentService(studentRepository)
	studentHandler := NewStudentHandler(studentService)

	prestasiRepository := repository.NewPrestasiRepository(pool)
	prestasiService := service.NewPrestasiService(studentRepository, prestasiRepository)
	prestasiHandler := NewPrestasiHandler(prestasiService)

	userRepository := repository.NewUserRepository(pool)
	tokenRepository := repository.NewTokenRepository(pool)
	authService := service.NewAuthService(
		userRepository, tokenRepository, jwtManager,
		time.Duration(config.GetEnvInt("JWT_REFRESH_TTL_DAYS", 7))*24*time.Hour,
	)
	authHandler := NewAuthHandler(authService)

	// 4. Aplikasi
	app := config.NewApp(logger, route.Dependencies{
		Pool: pool,
		JWT:  jwtManager,
		Students: route.StudentHandlers{
			List:             studentHandler.List,
			Get:              studentHandler.Get,
			Create:           studentHandler.Create,
			Replace:          studentHandler.Replace,
			Patch:            studentHandler.Patch,
			Delete:           studentHandler.Delete,
			GetPrestasiByNIM: prestasiHandler.GetByNIM,
		},
		Auth: route.AuthHandlers{
			Register: authHandler.Register,
			Login:    authHandler.Login,
			Refresh:  authHandler.Refresh,
			Logout:   authHandler.Logout,
			Me:       authHandler.Me,
		},
	})

	port := config.GetEnv("APP_PORT", "3000")

	go func() {
		if err := app.Listen(":" + port); err != nil {
			logger.Error("server berhenti", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()
	logger.Info("server berjalan", slog.String("port", port))

	// 5. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("sinyal berhenti diterima, menutup server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("gagal menutup server dengan rapi", slog.String("error", err.Error()))
	}
	logger.Info("server berhenti dengan rapi")
}