package models

import (
	"context"
	"database/sql"
	"time"
)

type Godown struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone,omitempty"`
	Notes     *string   `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GodownModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *GodownModel) Insert(ctx context.Context, godown *Godown) (int64, error) {
	query := `
		INSERT INTO godowns (name, phone, notes)
		VALUES (?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		godown.Name,
		godown.Phone,
		godown.Notes,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

const godownSelectCols = `id, name, phone, notes, created_at, updated_at`

func (m *GodownModel) GetByID(ctx context.Context, id int64) (*Godown, error) {
	query := `SELECT ` + godownSelectCols + ` FROM godowns WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanGodown(row)
}

func (m *GodownModel) GetAll(ctx context.Context) ([]Godown, error) {
	query := `SELECT ` + godownSelectCols + ` FROM godowns ORDER BY name ASC`
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
		SELECT ` + godownSelectCols + `
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
		SET name = ?, phone = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		godown.Name,
		godown.Phone,
		godown.Notes,
		godown.ID,
	)
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
	err := row.Scan(
		&godown.ID,
		&godown.Name,
		&godown.Phone,
		&godown.Notes,
		&godown.CreatedAt,
		&godown.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return godown, nil
}

func (m *GodownModel) scanGodownRow(rows *sql.Rows) (*Godown, error) {
	godown := &Godown{}
	err := rows.Scan(
		&godown.ID,
		&godown.Name,
		&godown.Phone,
		&godown.Notes,
		&godown.CreatedAt,
		&godown.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return godown, nil
}
