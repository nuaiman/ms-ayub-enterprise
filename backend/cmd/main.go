package main

import (
	"backend/internal/app"
	"backend/internal/bootstrap"
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/handlers"
	"backend/internal/models"
	"backend/internal/router"
	"log"
)

func main() {
	cfg := config.MustLoadConfig()

	dbPool := db.InitDB(cfg.DBPath, cfg.DBName, cfg.SchemaPath)
	defer db.CloseDB(dbPool)

	// =============================================
	// SEED (OPTIONAL)
	// Comment this line out to skip seeding.
	// =============================================
	db.SeedDB(dbPool, "internal/db/seed/seed.sql")

	a := app.Application{
		Config: cfg,
		DB:     dbPool,
		Models: models.NewModel(dbPool),
	}

	// NON SEED (CRITICAL)
	// =============================================
	bootstrap.BootstrapAdmin(&a)

	h := handlers.New(&a)

	r := router.RegisterRouter(&a, h)

	log.Fatal(a.Run(r))
}

// go run ./cmd/main.go --config=config/config.env
