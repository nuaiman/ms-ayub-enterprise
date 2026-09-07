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

	a := app.Application{
		Config: cfg,
		DB:     dbPool,
		Models: models.NewModel(dbPool),
	}

	// NON SEED (CRTICAL)
	// =============================================
	bootstrap.BootstrapAdmin(&a)

	// SEED DATABASE (Uncomment to run seed)
	// =============================================
	// log.Println("🌱 Seeding database...")

	// seedPath := filepath.Join("internal", "db", "schema", "seed.sql")

	// if _, err := os.Stat(seedPath); os.IsNotExist(err) {
	// 	log.Fatalf("❌ Seed file not found: %s", seedPath)
	// }

	// seedData, err := os.ReadFile(seedPath)
	// if err != nil {
	// 	log.Fatalf("❌ Failed to read seed file: %v", err)
	// }

	// _, err = dbPool.Exec(string(seedData))
	// if err != nil {
	// 	log.Fatalf("❌ Failed to execute seed: %v", err)
	// }

	// log.Println("✅ Seed data imported successfully!")
	// log.Println("📝 Default login credentials:")
	// log.Println("   Admin:    admin / password")
	// log.Println("   Manager:  manager1 / password")
	// log.Println("   Accounts: accounts1 / password")
	// log.Println("   Staff:    staff1 / password")
	// =============================================

	h := handlers.New(&a)

	r := router.RegisterRouter(&a, h)

	log.Fatal(a.Run(r))
}

// go run ./cmd/app/ --config=config/config.env
