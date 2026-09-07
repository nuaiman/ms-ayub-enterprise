package models

import (
	"context"
	"database/sql"
	"time"
)

type Godown struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Phone       *string   `json:"phone,omitempty"`
	Notes       *string   `json:"notes,omitempty"`
	IsActive    bool      `json:"is_active"`
	MonthlyRent float64   `json:"monthly_rent"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GodownModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *GodownModel) Insert(ctx context.Context, godown *Godown) (int64, error) {
	query := `
		INSERT INTO godowns (name, phone, notes, is_active, monthly_rent)
		VALUES (?, ?, ?, ?, ?)
	`

	isActive := 0
	if godown.IsActive {
		isActive = 1
	}

	res, err := m.DB.ExecContext(ctx, query,
		godown.Name,
		godown.Phone,
		godown.Notes,
		isActive,
		godown.MonthlyRent,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *GodownModel) GetByID(ctx context.Context, id int64) (*Godown, error) {
	query := `
		SELECT id, name, phone, notes, is_active, monthly_rent, created_at, updated_at
		FROM godowns
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanGodown(row)
}

func (m *GodownModel) GetAll(ctx context.Context) ([]Godown, error) {
	query := `
		SELECT id, name, phone, notes, is_active, monthly_rent, created_at, updated_at
		FROM godowns
		ORDER BY name ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	godowns := []Godown{}
	for rows.Next() {
		godown, err := m.scanGodownRow(rows)
		if err != nil {
			return nil, err
		}
		godowns = append(godowns, *godown)
	}

	return godowns, rows.Err()
}

func (m *GodownModel) GetActive(ctx context.Context) ([]Godown, error) {
	query := `
		SELECT id, name, phone, notes, is_active, monthly_rent, created_at, updated_at
		FROM godowns
		WHERE is_active = 1
		ORDER BY name ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	godowns := []Godown{}
	for rows.Next() {
		godown, err := m.scanGodownRow(rows)
		if err != nil {
			return nil, err
		}
		godowns = append(godowns, *godown)
	}

	return godowns, rows.Err()
}

func (m *GodownModel) Search(ctx context.Context, query string) ([]Godown, error) {
	searchQuery := `
		SELECT id, name, phone, notes, is_active, monthly_rent, created_at, updated_at
		FROM godowns
		WHERE name LIKE ? OR phone LIKE ?
		ORDER BY name ASC
	`

	searchTerm := "%" + query + "%"

	rows, err := m.DB.QueryContext(ctx, searchQuery, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	godowns := []Godown{}
	for rows.Next() {
		godown, err := m.scanGodownRow(rows)
		if err != nil {
			return nil, err
		}
		godowns = append(godowns, *godown)
	}

	return godowns, rows.Err()
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *GodownModel) Update(ctx context.Context, godown *Godown) error {
	query := `
		UPDATE godowns
		SET 
			name = ?,
			phone = ?,
			notes = ?,
			is_active = ?,
			monthly_rent = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	isActive := 0
	if godown.IsActive {
		isActive = 1
	}

	_, err := m.DB.ExecContext(ctx, query,
		godown.Name,
		godown.Phone,
		godown.Notes,
		isActive,
		godown.MonthlyRent,
		godown.ID,
	)
	return err
}

func (m *GodownModel) ToggleActive(ctx context.Context, id int64, isActive bool) error {
	active := 0
	if isActive {
		active = 1
	}

	query := `
		UPDATE godowns
		SET is_active = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, active, id)
	return err
}

func (m *GodownModel) UpdateMonthlyRent(ctx context.Context, id int64, monthlyRent float64) error {
	query := `
		UPDATE godowns
		SET monthly_rent = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, monthlyRent, id)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *GodownModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM godowns WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *GodownModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM godowns WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *GodownModel) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM godowns WHERE name = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, name).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *GodownModel) scanGodown(row *sql.Row) (*Godown, error) {
	godown := &Godown{}
	var isActive int
	err := row.Scan(
		&godown.ID,
		&godown.Name,
		&godown.Phone,
		&godown.Notes,
		&isActive,
		&godown.MonthlyRent,
		&godown.CreatedAt,
		&godown.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	godown.IsActive = isActive == 1
	return godown, nil
}

func (m *GodownModel) scanGodownRow(rows *sql.Rows) (*Godown, error) {
	godown := &Godown{}
	var isActive int
	err := rows.Scan(
		&godown.ID,
		&godown.Name,
		&godown.Phone,
		&godown.Notes,
		&isActive,
		&godown.MonthlyRent,
		&godown.CreatedAt,
		&godown.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	godown.IsActive = isActive == 1
	return godown, nil
}
