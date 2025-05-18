package main

import (
	"context"
	"fmt"
	_ "github.com/Gustcat/people-info-service/docs"
	"github.com/Gustcat/people-info-service/internal/config"
	"github.com/Gustcat/people-info-service/internal/http-server/handlers/persons"
	"github.com/Gustcat/people-info-service/internal/repository/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @title People info service API
// @version 1.0
// @description REST API-сервис для работы с информацией о людях.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("doesn't load env file: %s", err.Error())
		log.Fatalf("doesn't load env file: %s", err.Error())
	}

	conf, err := config.New()
	if err != nil {
		fmt.Printf("doesn't create config: %s", err.Error())
		log.Fatal(err)
	}
	fmt.Printf("addr: %s", conf.HTTPServer.Address)

	repo, err := postgres.NewRepo(ctx, conf.Postgres.DSN)
	if err != nil {
		fmt.Printf("doesn't create repo: %s", err.Error())
		log.Fatal(err)
	}
	defer repo.Close()

	router := chi.NewRouter()

	//router.Use(middleware.Recoverer)

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Route("/api/v1/persons", func(r chi.Router) {
		r.Post("/", persons.Create(ctx, repo))
		r.Get("/", persons.List(ctx, repo))
		r.Get("/{id}", persons.GetByID(ctx, repo))
		r.Patch("/{id}", persons.Update(ctx, repo))
		r.Delete("/{id}", persons.Delete(ctx, repo))
	})

	srv := &http.Server{
		Addr:         conf.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  conf.HTTPServer.Timeout,
		WriteTimeout: conf.HTTPServer.Timeout,
		IdleTimeout:  conf.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to start http server")
	}

}
