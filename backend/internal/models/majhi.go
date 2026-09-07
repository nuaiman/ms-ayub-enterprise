package models

import (
	"context"
	"database/sql"
	"time"
)

type Majhi struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone,omitempty"`
	Notes     *string   `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MajhiModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *MajhiModel) Insert(ctx context.Context, majhi *Majhi) (int64, error) {
	query := `
		INSERT INTO majhis (name, phone, notes)
		VALUES (?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		majhi.Name,
		majhi.Phone,
		majhi.Notes,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *MajhiModel) GetByID(ctx context.Context, id int64) (*Majhi, error) {
	query := `
		SELECT id, name, phone, notes, created_at, updated_at
		FROM majhis
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanMajhi(row)
}

func (m *MajhiModel) GetAll(ctx context.Context) ([]Majhi, error) {
	query := `
		SELECT id, name, phone, notes, created_at, updated_at
		FROM majhis
		ORDER BY name ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	majhis := []Majhi{}
	for rows.Next() {
		majhi, err := m.scanMajhiRow(rows)
		if err != nil {
			return nil, err
		}
		majhis = append(majhis, *majhi)
	}

	return majhis, rows.Err()
}

func (m *MajhiModel) Search(ctx context.Context, query string) ([]Majhi, error) {
	searchQuery := `
		SELECT id, name, phone, notes, created_at, updated_at
		FROM majhis
		WHERE name LIKE ? OR phone LIKE ?
		ORDER BY name ASC
	`

	searchTerm := "%" + query + "%"

	rows, err := m.DB.QueryContext(ctx, searchQuery, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	majhis := []Majhi{}
	for rows.Next() {
		majhi, err := m.scanMajhiRow(rows)
		if err != nil {
			return nil, err
		}
		majhis = append(majhis, *majhi)
	}

	return majhis, rows.Err()
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *MajhiModel) Update(ctx context.Context, majhi *Majhi) error {
	query := `
		UPDATE majhis
		SET name = ?, phone = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		majhi.Name,
		majhi.Phone,
		majhi.Notes,
		majhi.ID,
	)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *MajhiModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM majhis WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *MajhiModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM majhis WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *MajhiModel) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM majhis WHERE name = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, name).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *MajhiModel) scanMajhi(row *sql.Row) (*Majhi, error) {
	majhi := &Majhi{}
	err := row.Scan(
		&majhi.ID,
		&majhi.Name,
		&majhi.Phone,
		&majhi.Notes,
		&majhi.CreatedAt,
		&majhi.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return majhi, nil
}

func (m *MajhiModel) scanMajhiRow(rows *sql.Rows) (*Majhi, error) {
	majhi := &Majhi{}
	err := rows.Scan(
		&majhi.ID,
		&majhi.Name,
		&majhi.Phone,
		&majhi.Notes,
		&majhi.CreatedAt,
		&majhi.UpdatedAt,
	)
	return majhi, err
}
