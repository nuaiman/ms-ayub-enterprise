package bootstrap

import (
	"backend/internal/app"
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"log"
)

func BootstrapAdmin(app *app.Application) {
	if !app.Config.BootstrapAdmin {
		return
	}

	ctx := context.Background()

	existing, err := app.Models.User.GetByUsername(ctx, app.Config.BootstrapUsername)
	if err != nil {
		log.Fatalf("bootstrap check failed: %v", err)
	}

	if existing != nil {
		log.Println("bootstrap admin already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(app.Config.BootstrapPassword)
	if err != nil {
		log.Fatalf("failed to hash bootstrap password: %v", err)
	}

	admin := &models.User{
		Name:     app.Config.BootstrapAdminName,
		Username: app.Config.BootstrapUsername,
		Password: hashedPassword,
		Role:     "admin",
	}

	id, err := app.Models.User.Insert(ctx, admin)
	if err != nil {
		log.Fatalf("failed to create bootstrap admin: %v", err)
	}

	log.Printf("bootstrap admin created with id=%d", id)
}
