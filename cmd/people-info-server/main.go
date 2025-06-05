package main

import (
	"context"
	swaggerFiles "github.com/swaggo/files"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/Gustcat/people-info-service/docs"
	"github.com/Gustcat/people-info-service/internal/config"
	"github.com/Gustcat/people-info-service/internal/http-server/handlers/persons"
	"github.com/Gustcat/people-info-service/internal/logger"
	"github.com/Gustcat/people-info-service/internal/repository/postgres"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/gin-swagger"
)

const (
	envLocal = "local"
)

// @title People info service API
// @version 1.0
// @description REST API-сервис для работы с информацией о людях.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	log := logger.SetupLogger(slog.LevelInfo)

	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Warn("doesn't load env file: %s", slog.String("error", err.Error()))
	}

	conf, err := config.New()
	if err != nil {
		log.Error("doesn't set config: %s", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if conf.Env == envLocal {
		log = logger.SetupLogger(slog.LevelDebug)
	}

	log.Debug("Try to connect to db", slog.String("DSN", conf.Postgres.DSN))

	repo, err := postgres.NewRepo(ctx, conf.Postgres.DSN)
	if err != nil {
		log.Error("doesn't create repo", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer repo.Close()

	log.Debug("Try to setup router")
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r := router.Group("/api/v1/persons")
	{
		r.POST("/", persons.Create(log, repo))
		r.GET("/", persons.List(log, repo))
		r.GET("/:id", persons.GetByID(log, repo))
		r.PATCH("/:id", persons.Update(log, repo))
		r.DELETE("/:id", persons.Delete(log, repo))
	}

	srv := &http.Server{
		Addr:         conf.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  conf.HTTPServer.Timeout,
		WriteTimeout: conf.HTTPServer.Timeout,
		IdleTimeout:  conf.HTTPServer.IdleTimeout,
	}

	log.Info("Server started", slog.String("address", conf.HTTPServer.Address))

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start http server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
