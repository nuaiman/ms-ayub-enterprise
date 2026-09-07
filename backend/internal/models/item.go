package models

import (
	"context"
	"database/sql"
	"time"
)

type Item struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CustomerID  *int64    `json:"customer_id,omitempty"`
	ProductName *string   `json:"product_name,omitempty"`
	Category    *string   `json:"category,omitempty"`
	IsActive    bool      `json:"is_active"`
	Notes       *string   `json:"notes,omitempty"`
	ImageURL    *string   `json:"image_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ItemModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *ItemModel) Insert(ctx context.Context, item *Item) (int64, error) {
	query := `
		INSERT INTO items (user_id, customer_id, product_name, category, is_active, notes, image_url)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	isActive := 0
	if item.IsActive {
		isActive = 1
	}

	res, err := m.DB.ExecContext(ctx, query,
		item.UserID,
		item.CustomerID,
		item.ProductName,
		item.Category,
		isActive,
		item.Notes,
		item.ImageURL,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *ItemModel) GetByID(ctx context.Context, id int64) (*Item, error) {
	query := `
		SELECT id, user_id, customer_id, product_name, category, is_active, notes, image_url, created_at, updated_at
		FROM items
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanItem(row)
}

func (m *ItemModel) GetAll(ctx context.Context) ([]Item, error) {
	query := `
		SELECT id, user_id, customer_id, product_name, category, is_active, notes, image_url, created_at, updated_at
		FROM items
		ORDER BY product_name ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		item, err := m.scanItemRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, rows.Err()
}

func (m *ItemModel) GetByUserID(ctx context.Context, userID int64) ([]Item, error) {
	query := `
		SELECT id, user_id, customer_id, product_name, category, is_active, notes, image_url, created_at, updated_at
		FROM items
		WHERE user_id = ?
		ORDER BY product_name ASC
	`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		item, err := m.scanItemRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, rows.Err()
}

func (m *ItemModel) GetByCustomerID(ctx context.Context, customerID int64) ([]Item, error) {
	query := `
		SELECT id, user_id, customer_id, product_name, category, is_active, notes, image_url, created_at, updated_at
		FROM items
		WHERE customer_id = ?
		ORDER BY product_name ASC
	`

	rows, err := m.DB.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		item, err := m.scanItemRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, rows.Err()
}

func (m *ItemModel) GetActive(ctx context.Context) ([]Item, error) {
	query := `
		SELECT id, user_id, customer_id, product_name, category, is_active, notes, image_url, created_at, updated_at
		FROM items
		WHERE is_active = 1
		ORDER BY product_name ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		item, err := m.scanItemRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, rows.Err()
}

func (m *ItemModel) Search(ctx context.Context, query string) ([]Item, error) {
	searchQuery := `
		SELECT id, user_id, customer_id, product_name, category, is_active, notes, image_url, created_at, updated_at
		FROM items
		WHERE product_name LIKE ? OR category LIKE ?
		ORDER BY product_name ASC
	`

	searchTerm := "%" + query + "%"

	rows, err := m.DB.QueryContext(ctx, searchQuery, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		item, err := m.scanItemRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, rows.Err()
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *ItemModel) Update(ctx context.Context, item *Item) error {
	query := `
		UPDATE items
		SET 
			customer_id = ?,
			product_name = ?,
			category = ?,
			is_active = ?,
			notes = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	isActive := 0
	if item.IsActive {
		isActive = 1
	}

	_, err := m.DB.ExecContext(ctx, query,
		item.CustomerID,
		item.ProductName,
		item.Category,
		isActive,
		item.Notes,
		item.ID,
	)
	return err
}

func (m *ItemModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	query := `
		UPDATE items
		SET image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, imageURL, id)
	return err
}

func (m *ItemModel) ToggleActive(ctx context.Context, id int64, isActive bool) error {
	active := 0
	if isActive {
		active = 1
	}

	query := `
		UPDATE items
		SET is_active = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, active, id)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *ItemModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM items WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *ItemModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM items WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *ItemModel) CountByUser(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM items WHERE user_id = ?`

	var count int
	err := m.DB.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *ItemModel) scanItem(row *sql.Row) (*Item, error) {
	item := &Item{}
	err := row.Scan(
		&item.ID,
		&item.UserID,
		&item.CustomerID,
		&item.ProductName,
		&item.Category,
		&item.IsActive,
		&item.Notes,
		&item.ImageURL,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (m *ItemModel) scanItemRow(rows *sql.Rows) (*Item, error) {
	item := &Item{}
	err := rows.Scan(
		&item.ID,
		&item.UserID,
		&item.CustomerID,
		&item.ProductName,
		&item.Category,
		&item.IsActive,
		&item.Notes,
		&item.ImageURL,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	return item, err
}
