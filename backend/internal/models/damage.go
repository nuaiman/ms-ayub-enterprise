package models

import (
	"context"
	"database/sql"
	"time"
)

type Damage struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	StoreID      int64     `json:"store_id"`
	Quantity     float64   `json:"quantity"`
	QuantityUnit string    `json:"quantity_unit"`
	Weight       float64   `json:"weight"`
	WeightUnit   string    `json:"weight_unit"`
	DamageDate   time.Time `json:"damage_date"`
	Reason       string    `json:"reason"`
	Amount       float64   `json:"amount"`
	Notes        *string   `json:"notes,omitempty"`
	ImageURL     *string   `json:"image_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type DamageModel struct {
	DB *sql.DB
}

const damageSelectCols = `
	id, user_id, store_id,
	quantity, quantity_unit, weight, weight_unit,
	damage_date, reason, amount, notes, image_url,
	created_at, updated_at
`

// =============================================================================
// CREATE
// =============================================================================

func (m *DamageModel) Insert(ctx context.Context, d *Damage) (int64, error) {
	query := `
		INSERT INTO damages (
			user_id, store_id,
			quantity, quantity_unit, weight, weight_unit,
			damage_date, reason, amount, notes, image_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		d.UserID,
		d.StoreID,
		d.Quantity,
		d.QuantityUnit,
		d.Weight,
		d.WeightUnit,
		d.DamageDate,
		d.Reason,
		d.Amount,
		d.Notes,
		d.ImageURL,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *DamageModel) GetByID(ctx context.Context, id int64) (*Damage, error) {
	query := `SELECT ` + damageSelectCols + ` FROM damages WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scan(row)
}

func (m *DamageModel) GetAll(ctx context.Context) ([]Damage, error) {
	query := `SELECT ` + damageSelectCols + ` FROM damages ORDER BY damage_date DESC, id DESC`
	return m.query(ctx, query)
}

func (m *DamageModel) GetByStoreID(ctx context.Context, storeID int64) ([]Damage, error) {
	query := `SELECT ` + damageSelectCols + ` FROM damages WHERE store_id = ? ORDER BY damage_date DESC, id DESC`
	return m.query(ctx, query, storeID)
}

func (m *DamageModel) GetByUserID(ctx context.Context, userID int64) ([]Damage, error) {
	query := `SELECT ` + damageSelectCols + ` FROM damages WHERE user_id = ? ORDER BY damage_date DESC, id DESC`
	return m.query(ctx, query, userID)
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *DamageModel) Update(ctx context.Context, d *Damage) error {
	query := `
		UPDATE damages
		SET
			quantity = ?,
			quantity_unit = ?,
			weight = ?,
			weight_unit = ?,
			damage_date = ?,
			reason = ?,
			amount = ?,
			notes = ?,
			image_url = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := m.DB.ExecContext(ctx, query,
		d.Quantity,
		d.QuantityUnit,
		d.Weight,
		d.WeightUnit,
		d.DamageDate,
		d.Reason,
		d.Amount,
		d.Notes,
		d.ImageURL,
		d.ID,
	)
	return err
}

func (m *DamageModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	query := `
		UPDATE damages
		SET image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := m.DB.ExecContext(ctx, query, imageURL, id)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *DamageModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM damages WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *DamageModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM damages WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS / HELPERS
// =============================================================================

func (m *DamageModel) query(ctx context.Context, query string, args ...any) ([]Damage, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	damages := []Damage{}
	for rows.Next() {
		d, err := m.scanRow(rows)
		if err != nil {
			return nil, err
		}
		damages = append(damages, *d)
	}
	return damages, rows.Err()
}

func (m *DamageModel) scan(row *sql.Row) (*Damage, error) {
	d := &Damage{}
	err := row.Scan(
		&d.ID,
		&d.UserID,
		&d.StoreID,
		&d.Quantity,
		&d.QuantityUnit,
		&d.Weight,
		&d.WeightUnit,
		&d.DamageDate,
		&d.Reason,
		&d.Amount,
		&d.Notes,
		&d.ImageURL,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return d, nil
}

func (m *DamageModel) scanRow(rows *sql.Rows) (*Damage, error) {
	d := &Damage{}
	err := rows.Scan(
		&d.ID,
		&d.UserID,
		&d.StoreID,
		&d.Quantity,
		&d.QuantityUnit,
		&d.Weight,
		&d.WeightUnit,
		&d.DamageDate,
		&d.Reason,
		&d.Amount,
		&d.Notes,
		&d.ImageURL,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	return d, err
}
