package models

import (
	"context"
	"database/sql"
	"time"
)

type StoreAdjustment struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	StoreID        int64     `json:"store_id"`
	AdjustmentType string    `json:"adjustment_type"` // "delta" | "absolute"
	InputWeight    float64   `json:"input_weight"`
	InputQuantity  float64   `json:"input_quantity"`
	WeightDelta    float64   `json:"weight_delta"`
	QuantityDelta  float64   `json:"quantity_delta"`
	Reason         *string   `json:"reason,omitempty"`
	Notes          *string   `json:"notes,omitempty"`
	AdjustedAt     time.Time `json:"adjusted_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type StoreAdjustmentModel struct {
	DB *sql.DB
}

const storeAdjustmentSelectCols = `
	id, user_id, store_id,
	adjustment_type,
	input_weight, input_quantity,
	weight_delta, quantity_delta,
	reason, notes,
	adjusted_at, created_at, updated_at
`

// =============================================================================
// CREATE
// =============================================================================

func (m *StoreAdjustmentModel) Insert(ctx context.Context, adj *StoreAdjustment) (int64, error) {
	query := `
		INSERT INTO store_adjustments (
			user_id, store_id,
			adjustment_type,
			input_weight, input_quantity,
			weight_delta, quantity_delta,
			reason, notes,
			adjusted_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		adj.UserID,
		adj.StoreID,
		adj.AdjustmentType,
		adj.InputWeight,
		adj.InputQuantity,
		adj.WeightDelta,
		adj.QuantityDelta,
		adj.Reason,
		adj.Notes,
		adj.AdjustedAt,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *StoreAdjustmentModel) GetByID(ctx context.Context, id int64) (*StoreAdjustment, error) {
	query := `SELECT ` + storeAdjustmentSelectCols + ` FROM store_adjustments WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scan(row)
}

func (m *StoreAdjustmentModel) GetByStoreID(ctx context.Context, storeID int64) ([]StoreAdjustment, error) {
	query := `SELECT ` + storeAdjustmentSelectCols + `
		FROM store_adjustments
		WHERE store_id = ?
		ORDER BY adjusted_at DESC, id DESC`
	return m.query(ctx, query, storeID)
}

func (m *StoreAdjustmentModel) GetAll(ctx context.Context) ([]StoreAdjustment, error) {
	query := `SELECT ` + storeAdjustmentSelectCols + `
		FROM store_adjustments
		ORDER BY adjusted_at DESC, id DESC`
	return m.query(ctx, query)
}

func (m *StoreAdjustmentModel) GetByUserID(ctx context.Context, userID int64) ([]StoreAdjustment, error) {
	query := `SELECT ` + storeAdjustmentSelectCols + `
		FROM store_adjustments
		WHERE user_id = ?
		ORDER BY adjusted_at DESC, id DESC`
	return m.query(ctx, query, userID)
}

// =============================================================================
// DELETE
// =============================================================================

func (m *StoreAdjustmentModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM store_adjustments WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *StoreAdjustmentModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM store_adjustments WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS / HELPERS
// =============================================================================

func (m *StoreAdjustmentModel) query(ctx context.Context, query string, args ...any) ([]StoreAdjustment, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	adjustments := []StoreAdjustment{}
	for rows.Next() {
		a, err := m.scanRow(rows)
		if err != nil {
			return nil, err
		}
		adjustments = append(adjustments, *a)
	}
	return adjustments, rows.Err()
}

func (m *StoreAdjustmentModel) scan(row *sql.Row) (*StoreAdjustment, error) {
	a := &StoreAdjustment{}
	err := row.Scan(
		&a.ID,
		&a.UserID,
		&a.StoreID,
		&a.AdjustmentType,
		&a.InputWeight,
		&a.InputQuantity,
		&a.WeightDelta,
		&a.QuantityDelta,
		&a.Reason,
		&a.Notes,
		&a.AdjustedAt,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}

func (m *StoreAdjustmentModel) scanRow(rows *sql.Rows) (*StoreAdjustment, error) {
	a := &StoreAdjustment{}
	err := rows.Scan(
		&a.ID,
		&a.UserID,
		&a.StoreID,
		&a.AdjustmentType,
		&a.InputWeight,
		&a.InputQuantity,
		&a.WeightDelta,
		&a.QuantityDelta,
		&a.Reason,
		&a.Notes,
		&a.AdjustedAt,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	return a, err
}
