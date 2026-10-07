package models

import (
	"context"
	"database/sql"
	"time"
)

type StoreTransfer struct {
	ID                 int64     `json:"id"`
	UserID             int64     `json:"user_id"`
	StoreID            int64     `json:"store_id"`
	FromGodownID       int64     `json:"from_godown_id"`
	ToGodownID         int64     `json:"to_godown_id"`
	WeightAtTransfer   float64   `json:"weight_at_transfer"`
	QuantityAtTransfer float64   `json:"quantity_at_transfer"`
	Notes              *string   `json:"notes,omitempty"`
	TransferredAt      time.Time `json:"transferred_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type StoreTransferModel struct {
	DB *sql.DB
}

const storeTransferSelectCols = `
	id, user_id, store_id,
	from_godown_id, to_godown_id,
	weight_at_transfer, quantity_at_transfer,
	notes,
	transferred_at, created_at, updated_at
`

// =============================================================================
// CREATE
// =============================================================================

func (m *StoreTransferModel) Insert(ctx context.Context, t *StoreTransfer) (int64, error) {
	query := `
		INSERT INTO store_transfers (
			user_id, store_id,
			from_godown_id, to_godown_id,
			weight_at_transfer, quantity_at_transfer,
			notes,
			transferred_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		t.UserID,
		t.StoreID,
		t.FromGodownID,
		t.ToGodownID,
		t.WeightAtTransfer,
		t.QuantityAtTransfer,
		t.Notes,
		t.TransferredAt,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *StoreTransferModel) GetByID(ctx context.Context, id int64) (*StoreTransfer, error) {
	query := `SELECT ` + storeTransferSelectCols + ` FROM store_transfers WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scan(row)
}

func (m *StoreTransferModel) GetByStoreID(ctx context.Context, storeID int64) ([]StoreTransfer, error) {
	query := `SELECT ` + storeTransferSelectCols + `
		FROM store_transfers
		WHERE store_id = ?
		ORDER BY transferred_at DESC, id DESC`
	return m.query(ctx, query, storeID)
}

func (m *StoreTransferModel) GetAll(ctx context.Context) ([]StoreTransfer, error) {
	query := `SELECT ` + storeTransferSelectCols + `
		FROM store_transfers
		ORDER BY transferred_at DESC, id DESC`
	return m.query(ctx, query)
}

func (m *StoreTransferModel) GetByUserID(ctx context.Context, userID int64) ([]StoreTransfer, error) {
	query := `SELECT ` + storeTransferSelectCols + `
		FROM store_transfers
		WHERE user_id = ?
		ORDER BY transferred_at DESC, id DESC`
	return m.query(ctx, query, userID)
}

// GetLatestForStore returns the most recent transfer for a store, or nil if none.
func (m *StoreTransferModel) GetLatestForStore(ctx context.Context, storeID int64) (*StoreTransfer, error) {
	query := `SELECT ` + storeTransferSelectCols + `
		FROM store_transfers
		WHERE store_id = ?
		ORDER BY transferred_at DESC, id DESC
		LIMIT 1`
	row := m.DB.QueryRowContext(ctx, query, storeID)
	return m.scan(row)
}

// =============================================================================
// DELETE
// =============================================================================

func (m *StoreTransferModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM store_transfers WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *StoreTransferModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM store_transfers WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS / HELPERS
// =============================================================================

func (m *StoreTransferModel) query(ctx context.Context, query string, args ...any) ([]StoreTransfer, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transfers := []StoreTransfer{}
	for rows.Next() {
		t, err := m.scanRow(rows)
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, *t)
	}
	return transfers, rows.Err()
}

func (m *StoreTransferModel) scan(row *sql.Row) (*StoreTransfer, error) {
	t := &StoreTransfer{}
	err := row.Scan(
		&t.ID,
		&t.UserID,
		&t.StoreID,
		&t.FromGodownID,
		&t.ToGodownID,
		&t.WeightAtTransfer,
		&t.QuantityAtTransfer,
		&t.Notes,
		&t.TransferredAt,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}

func (m *StoreTransferModel) scanRow(rows *sql.Rows) (*StoreTransfer, error) {
	t := &StoreTransfer{}
	err := rows.Scan(
		&t.ID,
		&t.UserID,
		&t.StoreID,
		&t.FromGodownID,
		&t.ToGodownID,
		&t.WeightAtTransfer,
		&t.QuantityAtTransfer,
		&t.Notes,
		&t.TransferredAt,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	return t, err
}
