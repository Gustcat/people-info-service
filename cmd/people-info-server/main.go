package main

import (
	"context"
	_ "github.com/Gustcat/people-info-service/docs"
	"github.com/Gustcat/people-info-service/internal/config"
	"github.com/Gustcat/people-info-service/internal/http-server/handlers/persons"
	"github.com/Gustcat/people-info-service/internal/logger"
	"github.com/Gustcat/people-info-service/internal/repository/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"os"
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
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Route("/api/v1/persons", func(r chi.Router) {
		r.Post("/", persons.Create(ctx, log, repo))
		r.Get("/", persons.List(ctx, log, repo))
		r.Get("/{id}", persons.GetByID(ctx, log, repo))
		r.Patch("/{id}", persons.Update(ctx, log, repo))
		r.Delete("/{id}", persons.Delete(ctx, log, repo))
	})

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
