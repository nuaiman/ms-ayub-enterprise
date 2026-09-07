package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (h *Handler) BackupHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[BACKUP] BackupHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	timestamp := time.Now().Format("2006-01-02_150405")
	filename := fmt.Sprintf("backup_%s.zip", timestamp)

	log.Printf("[BACKUP] Starting backup: %s", filename)

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	// backup data directory
	log.Printf("[BACKUP] Backing up data directory...")
	if err := addDirectoryToZip(zipWriter, "./data", "data"); err != nil {
		log.Printf("[BACKUP ERROR] Failed to backup data directory: %v", err)
		// Cannot send JSON error - headers already set for zip download
		return
	}

	// backup bucket directory
	log.Printf("[BACKUP] Backing up bucket directory...")
	if err := addDirectoryToZip(zipWriter, "./bucket", "bucket"); err != nil {
		log.Printf("[BACKUP ERROR] Failed to backup bucket directory: %v", err)
		return
	}

	log.Printf("[BACKUP] Backup completed: %s", filename)
}

func addDirectoryToZip(zipWriter *zip.Writer, sourceDir string, zipRoot string) error {
	// Check if directory exists
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		log.Printf("[BACKUP] Directory does not exist, skipping: %s", sourceDir)
		return nil
	}

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("[BACKUP] Error accessing %s: %v", path, err)
			return nil // Skip problematic files
		}

		if info.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			log.Printf("[BACKUP] Error getting relative path: %v", err)
			return nil
		}

		// Use forward slashes for zip compatibility
		zipPath := filepath.ToSlash(filepath.Join(zipRoot, relativePath))

		return addFileToZip(zipWriter, path, zipPath)
	})
}

func addFileToZip(zipWriter *zip.Writer, filePath string, zipName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(zipName)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}
