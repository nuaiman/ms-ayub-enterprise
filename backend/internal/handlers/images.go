// backend/internal/handlers/images.go
package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"backend/internal/middlewares"
	"backend/internal/utils"
)

func (h *Handler) UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	entityType := utils.GetURLParam(r, "type")

	id, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid id")
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid form")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		utils.ErrorJson(w, http.StatusBadRequest, "image required")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := "img_" + strconv.FormatInt(id, 10) + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + ext

	var (
		basePath     string
		imageURL     string
		oldImagePath string
	)

	currentUser, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil || currentUser == nil {
		log.Printf("[IMAGES] ERROR: failed to fetch current user - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify permissions")
		return
	}

	switch entityType {

	case "users":
		user, err := h.app.Models.User.GetByID(r.Context(), id)
		if err != nil {
			log.Printf("[IMAGES] ERROR: Error fetching user: %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
			return
		}
		if user == nil {
			utils.ErrorJson(w, http.StatusNotFound, "user not found")
			return
		}

		if user.ID != userID {
			if currentUser.Role != "admin" && currentUser.Role != "manager" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if user.ImageURL != nil && *user.ImageURL != "" {
			oldImagePath = "." + *user.ImageURL
		}

		basePath = filepath.Join("bucket", "users", strconv.FormatInt(id, 10))

	case "stores":
		store, err := h.app.Models.Store.GetByID(r.Context(), id)
		if err != nil {
			log.Printf("[IMAGES] ERROR: Error fetching store: %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch store")
			return
		}
		if store == nil {
			utils.ErrorJson(w, http.StatusNotFound, "store not found")
			return
		}

		if currentUser.Role != "admin" && currentUser.Role != "manager" {
			utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
			return
		}

		if store.ImageURL != nil && *store.ImageURL != "" {
			oldImagePath = "." + *store.ImageURL
		}

		basePath = filepath.Join("bucket", "stores", strconv.FormatInt(id, 10))

	case "expenses":
		expense, err := h.app.Models.Expense.GetByID(r.Context(), id)
		if err != nil {
			log.Printf("[IMAGES] ERROR: Error fetching expense: %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch expense")
			return
		}
		if expense == nil {
			utils.ErrorJson(w, http.StatusNotFound, "expense not found")
			return
		}

		if expense.UserID != userID {
			if currentUser.Role != "admin" && currentUser.Role != "manager" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if expense.ImageURL != nil && *expense.ImageURL != "" {
			oldImagePath = "." + *expense.ImageURL
		}

		basePath = filepath.Join("bucket", "expenses", strconv.FormatInt(id, 10))

	case "incomes":
		income, err := h.app.Models.Income.GetByID(r.Context(), id)
		if err != nil {
			log.Printf("[IMAGES] ERROR: Error fetching income: %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch income")
			return
		}
		if income == nil {
			utils.ErrorJson(w, http.StatusNotFound, "income not found")
			return
		}

		if income.UserID != userID {
			if currentUser.Role != "admin" && currentUser.Role != "manager" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if income.ImageURL != nil && *income.ImageURL != "" {
			oldImagePath = "." + *income.ImageURL
		}

		basePath = filepath.Join("bucket", "incomes", strconv.FormatInt(id, 10))

	case "deliveries":
		delivery, err := h.app.Models.Delivery.GetByID(r.Context(), id)
		if err != nil {
			log.Printf("[IMAGES] ERROR: Error fetching delivery: %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch delivery")
			return
		}
		if delivery == nil {
			utils.ErrorJson(w, http.StatusNotFound, "delivery not found")
			return
		}

		if delivery.UserID != userID {
			if currentUser.Role != "admin" && currentUser.Role != "manager" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if delivery.ImageURL != nil && *delivery.ImageURL != "" {
			oldImagePath = "." + *delivery.ImageURL
		}

		basePath = filepath.Join("bucket", "deliveries", strconv.FormatInt(id, 10))

	case "damages":
		damage, err := h.app.Models.Damage.GetByID(r.Context(), id)
		if err != nil {
			log.Printf("[IMAGES] ERROR: Error fetching damage: %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch damage")
			return
		}
		if damage == nil {
			utils.ErrorJson(w, http.StatusNotFound, "damage not found")
			return
		}

		if damage.UserID != userID {
			if currentUser.Role != "admin" && currentUser.Role != "manager" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if damage.ImageURL != nil && *damage.ImageURL != "" {
			oldImagePath = "." + *damage.ImageURL
		}

		basePath = filepath.Join("bucket", "damages", strconv.FormatInt(id, 10))

	default:
		utils.ErrorJson(w, http.StatusBadRequest, "invalid type. Supported: users, stores, expenses, incomes, deliveries, damages")
		return
	}

	err = os.MkdirAll(basePath, 0755)
	if err != nil {
		log.Printf("[IMAGES] ERROR: Failed to create directory: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to create directory")
		return
	}

	dstPath := filepath.Join(basePath, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		log.Printf("[IMAGES] ERROR: Failed to create file: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("[IMAGES] ERROR: Failed to write file: %v", err)
		os.Remove(dstPath)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to write file")
		return
	}

	imageURL = "/" + filepath.ToSlash(dstPath)

	switch entityType {
	case "users":
		err = h.app.Models.User.UpdateImage(r.Context(), id, &imageURL)
	case "stores":
		err = h.app.Models.Store.UpdateImage(r.Context(), id, &imageURL)
	case "expenses":
		err = h.app.Models.Expense.UpdateImage(r.Context(), id, &imageURL)
	case "incomes":
		err = h.app.Models.Income.UpdateImage(r.Context(), id, &imageURL)
	case "deliveries":
		err = h.app.Models.Delivery.UpdateImage(r.Context(), id, &imageURL)
	case "damages":
		err = h.app.Models.Damage.UpdateImage(r.Context(), id, &imageURL)
	}

	if err != nil {
		log.Printf("[IMAGES] ERROR: Failed to update database: %v", err)
		os.Remove(dstPath)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update database")
		return
	}

	if oldImagePath != "" {
		if err := os.Remove(oldImagePath); err != nil {
			log.Printf("[IMAGES] WARN: Failed to delete old image: %s - %v", oldImagePath, err)
		}
	}

	log.Printf("[IMAGES] SUCCESS: Image uploaded for %s ID=%d: %s", entityType, id, imageURL)
	utils.SuccessJson(w, http.StatusOK, "image uploaded", map[string]any{
		"image_url": imageURL,
	})
}

// DeleteImageHandler - DELETE /api/images/{type}/{id}
func (h *Handler) DeleteImageHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(int64)
	if !ok {
		utils.ErrorJson(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	entityType := utils.GetURLParam(r, "type")

	id, ok := utils.GetParamID(w, r)
	if !ok {
		utils.ErrorJson(w, http.StatusBadRequest, "invalid id")
		return
	}

	var imageURL string

	currentUser, err := h.app.Models.User.GetByID(r.Context(), userID)
	if err != nil || currentUser == nil {
		log.Printf("[IMAGES] ERROR: failed to fetch current user - %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to verify permissions")
		return
	}

	switch entityType {

	case "users":
		user, err := h.app.Models.User.GetByID(r.Context(), id)
		if err != nil {
			log.Printf("[IMAGES] ERROR: Error fetching user: %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch user")
			return
		}
		if user == nil {
			utils.ErrorJson(w, http.StatusNotFound, "user not found")
			return
		}

		if user.ID != userID {
			if currentUser.Role != "admin" && currentUser.Role != "manager" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if user.ImageURL != nil {
			imageURL = *user.ImageURL
		}

	case "stores":
		store, err := h.app.Models.Store.GetByID(r.Context(), id)
		if err != nil {
			log.Printf("[IMAGES] ERROR: Error fetching store: %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch store")
			return
		}
		if store == nil {
			utils.ErrorJson(w, http.StatusNotFound, "store not found")
			return
		}

		if currentUser.Role != "admin" && currentUser.Role != "manager" {
			utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
			return
		}

		if store.ImageURL != nil {
			imageURL = *store.ImageURL
		}

	case "expenses":
		expense, err := h.app.Models.Expense.GetByID(r.Context(), id)
		if err != nil {
			log.Printf("[IMAGES] ERROR: Error fetching expense: %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch expense")
			return
		}
		if expense == nil {
			utils.ErrorJson(w, http.StatusNotFound, "expense not found")
			return
		}

		if expense.UserID != userID {
			if currentUser.Role != "admin" && currentUser.Role != "manager" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if expense.ImageURL != nil {
			imageURL = *expense.ImageURL
		}

	case "incomes":
		income, err := h.app.Models.Income.GetByID(r.Context(), id)
		if err != nil {
			log.Printf("[IMAGES] ERROR: Error fetching income: %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch income")
			return
		}
		if income == nil {
			utils.ErrorJson(w, http.StatusNotFound, "income not found")
			return
		}

		if income.UserID != userID {
			if currentUser.Role != "admin" && currentUser.Role != "manager" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if income.ImageURL != nil {
			imageURL = *income.ImageURL
		}

	case "deliveries":
		delivery, err := h.app.Models.Delivery.GetByID(r.Context(), id)
		if err != nil {
			log.Printf("[IMAGES] ERROR: Error fetching delivery: %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch delivery")
			return
		}
		if delivery == nil {
			utils.ErrorJson(w, http.StatusNotFound, "delivery not found")
			return
		}

		if delivery.UserID != userID {
			if currentUser.Role != "admin" && currentUser.Role != "manager" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if delivery.ImageURL != nil {
			imageURL = *delivery.ImageURL
		}

	case "damages":
		damage, err := h.app.Models.Damage.GetByID(r.Context(), id)
		if err != nil {
			log.Printf("[IMAGES] ERROR: Error fetching damage: %v", err)
			utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch damage")
			return
		}
		if damage == nil {
			utils.ErrorJson(w, http.StatusNotFound, "damage not found")
			return
		}

		if damage.UserID != userID {
			if currentUser.Role != "admin" && currentUser.Role != "manager" {
				utils.ErrorJson(w, http.StatusForbidden, "unauthorized")
				return
			}
		}

		if damage.ImageURL != nil {
			imageURL = *damage.ImageURL
		}

	default:
		utils.ErrorJson(w, http.StatusBadRequest, "invalid type. Supported: users, stores, expenses, incomes, deliveries, damages")
		return
	}

	if imageURL == "" {
		utils.SuccessJson(w, http.StatusOK, "no image to delete", nil)
		return
	}

	if err := os.Remove("." + imageURL); err != nil {
		log.Printf("[IMAGES] WARN: Failed to delete image file: %s - %v", imageURL, err)
	}

	switch entityType {
	case "users":
		err = h.app.Models.User.UpdateImage(r.Context(), id, nil)
	case "stores":
		err = h.app.Models.Store.UpdateImage(r.Context(), id, nil)
	case "expenses":
		err = h.app.Models.Expense.UpdateImage(r.Context(), id, nil)
	case "incomes":
		err = h.app.Models.Income.UpdateImage(r.Context(), id, nil)
	case "deliveries":
		err = h.app.Models.Delivery.UpdateImage(r.Context(), id, nil)
	case "damages":
		err = h.app.Models.Damage.UpdateImage(r.Context(), id, nil)
	}

	if err != nil {
		log.Printf("[IMAGES] ERROR: Failed to update database: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to update database")
		return
	}

	log.Printf("[IMAGES] SUCCESS: Image deleted for %s ID=%d", entityType, id)
	utils.SuccessJson(w, http.StatusOK, "image deleted", nil)
}
