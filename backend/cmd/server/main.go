// Command server runs the Travel CRM REST API.
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"travelcrm/internal/config"
	"travelcrm/internal/db"
	"travelcrm/internal/handler"
	"travelcrm/internal/repository"
	"travelcrm/internal/service"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	repo := repository.New(pool)
	svc := service.New(repo)
	h := handler.New(svc)
	router := handler.Routes(h, cfg)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("Travel CRM API listening on http://localhost:%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server: %v", err)
	}
}
