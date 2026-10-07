package models

import (
	"context"
	"database/sql"
	"time"
)

type LotTransfer struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	LotID          int64     `json:"lot_id"`
	FromCustomerID int64     `json:"from_customer_id"`
	ToCustomerID   int64     `json:"to_customer_id"`
	Notes          *string   `json:"notes,omitempty"`
	TransferredAt  time.Time `json:"transferred_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type LotTransferModel struct {
	DB *sql.DB
}

const lotTransferSelectCols = `
	id, user_id, lot_id,
	from_customer_id, to_customer_id,
	notes,
	transferred_at, created_at, updated_at
`

// =============================================================================
// CREATE
// =============================================================================

func (m *LotTransferModel) Insert(ctx context.Context, t *LotTransfer) (int64, error) {
	query := `
		INSERT INTO lot_transfers (
			user_id, lot_id,
			from_customer_id, to_customer_id,
			notes,
			transferred_at
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		t.UserID,
		t.LotID,
		t.FromCustomerID,
		t.ToCustomerID,
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

func (m *LotTransferModel) GetByID(ctx context.Context, id int64) (*LotTransfer, error) {
	query := `SELECT ` + lotTransferSelectCols + ` FROM lot_transfers WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scan(row)
}

func (m *LotTransferModel) GetByLotID(ctx context.Context, lotID int64) ([]LotTransfer, error) {
	query := `SELECT ` + lotTransferSelectCols + `
		FROM lot_transfers
		WHERE lot_id = ?
		ORDER BY transferred_at DESC, id DESC`
	return m.query(ctx, query, lotID)
}

func (m *LotTransferModel) GetAll(ctx context.Context) ([]LotTransfer, error) {
	query := `SELECT ` + lotTransferSelectCols + `
		FROM lot_transfers
		ORDER BY transferred_at DESC, id DESC`
	return m.query(ctx, query)
}

func (m *LotTransferModel) GetByUserID(ctx context.Context, userID int64) ([]LotTransfer, error) {
	query := `SELECT ` + lotTransferSelectCols + `
		FROM lot_transfers
		WHERE user_id = ?
		ORDER BY transferred_at DESC, id DESC`
	return m.query(ctx, query, userID)
}

// GetLatestForLot returns the most recent transfer for a lot, or nil if none.
func (m *LotTransferModel) GetLatestForLot(ctx context.Context, lotID int64) (*LotTransfer, error) {
	query := `SELECT ` + lotTransferSelectCols + `
		FROM lot_transfers
		WHERE lot_id = ?
		ORDER BY transferred_at DESC, id DESC
		LIMIT 1`
	row := m.DB.QueryRowContext(ctx, query, lotID)
	return m.scan(row)
}

// =============================================================================
// DELETE
// =============================================================================

func (m *LotTransferModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM lot_transfers WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *LotTransferModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM lot_transfers WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS / HELPERS
// =============================================================================

func (m *LotTransferModel) query(ctx context.Context, query string, args ...any) ([]LotTransfer, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transfers := []LotTransfer{}
	for rows.Next() {
		t, err := m.scanRow(rows)
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, *t)
	}
	return transfers, rows.Err()
}

func (m *LotTransferModel) scan(row *sql.Row) (*LotTransfer, error) {
	t := &LotTransfer{}
	err := row.Scan(
		&t.ID,
		&t.UserID,
		&t.LotID,
		&t.FromCustomerID,
		&t.ToCustomerID,
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

func (m *LotTransferModel) scanRow(rows *sql.Rows) (*LotTransfer, error) {
	t := &LotTransfer{}
	err := rows.Scan(
		&t.ID,
		&t.UserID,
		&t.LotID,
		&t.FromCustomerID,
		&t.ToCustomerID,
		&t.Notes,
		&t.TransferredAt,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	return t, err
}
