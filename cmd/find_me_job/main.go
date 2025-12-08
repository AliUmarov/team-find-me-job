package main

import (
	"log/slog"
	"os"

	"github.com/AliUmarov/team-find-me-job/internal/config"
	"github.com/AliUmarov/team-find-me-job/internal/logger"
	"github.com/AliUmarov/team-find-me-job/internal/models"
	"github.com/AliUmarov/team-find-me-job/internal/repository"
	"github.com/AliUmarov/team-find-me-job/internal/services"
	"github.com/AliUmarov/team-find-me-job/internal/transport"
	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.InitLogger()

	config.SetEnv(log)
	db := config.Connect(log)

	if err := db.AutoMigrate(&models.Company{}); err != nil {
		log.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	log.Info("migrations completed")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	companyRepo := repository.NewCompanyRepository(db)

	companyService := services.NewCompanyService(*companyRepo)

	r := gin.Default()

	transport.RegisterRoutes(r, *companyService)

	log.Info("server started",
		slog.String("addr", port))

	if err := r.Run(":" + port); err != nil {
		log.Error("не удалось запустить сервер", slog.Any("error", err))
	}
}
