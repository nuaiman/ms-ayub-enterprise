package app

import (
	"backend/internal/config"
	"backend/internal/models"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	Config *config.Config
	DB     *sql.DB
	Models models.Models
}

func (app *Application) Run(r http.Handler) error {
	srv := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      r,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	errorCh := make(chan error, 1)

	go func() {
		log.Printf("Server is running at %s", app.Config.Addr)
		errorCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errorCh:
		return err

	case sig := <-signalCh:
		log.Printf("Shutdown signal recieved: %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("[X] graceful shutdown failed: %v", err)
			return err
		}

		app.DB.Close()

		log.Println("server closed gracefully")
	}

	return nil
}
