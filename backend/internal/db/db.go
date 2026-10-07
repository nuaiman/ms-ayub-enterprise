package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB(dbPath, dbName, schemaFile string) *sql.DB {
	if err := os.MkdirAll(dbPath, os.ModePerm); err != nil {
		log.Fatalf("failed to create database directory: %v", err)
	}

	dbFile := filepath.Join(dbPath, dbName)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	pragmas := []string{
		`PRAGMA foreign_keys = ON;`,
		`PRAGMA journal_mode = WAL;`,
		`PRAGMA synchronous = NORMAL;`,
		`PRAGMA busy_timeout = 5000;`,
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			log.Fatalf("failed to execute pragma: %v", err)
		}
	}

	if err := execSchema(db, schemaFile); err != nil {
		log.Fatalf("failed to execute schema (%s): %v", schemaFile, err)
	}

	log.Println("database initialized successfully")

	return db
}

// SeedDB applies a seed.sql file. Idempotent if the seed uses INSERT OR IGNORE.
// Missing seed file is a warning, not fatal — so production can skip it.
func SeedDB(db *sql.DB, seedFile string) {
	if _, err := os.Stat(seedFile); os.IsNotExist(err) {
		log.Printf("[SEED] no seed file at %s, skipping", seedFile)
		return
	}

	data, err := os.ReadFile(seedFile)
	if err != nil {
		log.Printf("[SEED] failed to read %s: %v", seedFile, err)
		return
	}

	if _, err := db.Exec(string(data)); err != nil {
		log.Printf("[SEED] failed to apply %s: %v", seedFile, err)
		return
	}

	log.Printf("[SEED] seed applied from %s", seedFile)
}

func execSchema(db *sql.DB, schemaFile string) error {
	data, err := os.ReadFile(schemaFile)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(data))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

func CloseDB(db *sql.DB) {
	if db == nil {
		return
	}

	if err := db.Close(); err != nil {
		log.Printf("error closing database: %v", err)
		return
	}

	log.Println("database closed")
}
