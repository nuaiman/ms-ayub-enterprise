package models

import (
	"context"
	"database/sql"
	"time"
)

type Discount struct {
	ID             int64     `json:"id"`
	InvoiceID      int64     `json:"invoice_id"`
	Type           string    `json:"type"`
	Value          float64   `json:"value"`
	ComputedAmount float64   `json:"computed_amount"`
	Reason         *string   `json:"reason,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type DiscountModel struct {
	DB *sql.DB
}

const discountSelectCols = `
	id, invoice_id, type, value, computed_amount,
	reason,
	created_at, updated_at
`

// =============================================================================
// CREATE
// =============================================================================

func (m *DiscountModel) Insert(ctx context.Context, d *Discount) (int64, error) {
	query := `
		INSERT INTO discounts (
			invoice_id, type, value, computed_amount, reason
		) VALUES (?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		d.InvoiceID,
		d.Type,
		d.Value,
		d.ComputedAmount,
		d.Reason,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *DiscountModel) GetByID(ctx context.Context, id int64) (*Discount, error) {
	query := `SELECT ` + discountSelectCols + ` FROM discounts WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scan(row)
}

func (m *DiscountModel) GetByInvoiceID(ctx context.Context, invoiceID int64) ([]Discount, error) {
	query := `SELECT ` + discountSelectCols + `
		FROM discounts
		WHERE invoice_id = ?
		ORDER BY id ASC`
	return m.query(ctx, query, invoiceID)
}

// =============================================================================
// DELETE
// =============================================================================

func (m *DiscountModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM discounts WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m *DiscountModel) DeleteByInvoiceID(ctx context.Context, invoiceID int64) error {
	query := `DELETE FROM discounts WHERE invoice_id = ?`
	_, err := m.DB.ExecContext(ctx, query, invoiceID)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *DiscountModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM discounts WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS / HELPERS
// =============================================================================

func (m *DiscountModel) query(ctx context.Context, query string, args ...any) ([]Discount, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	discounts := []Discount{}
	for rows.Next() {
		d, err := m.scanRow(rows)
		if err != nil {
			return nil, err
		}
		discounts = append(discounts, *d)
	}
	return discounts, rows.Err()
}

func (m *DiscountModel) scan(row *sql.Row) (*Discount, error) {
	d := &Discount{}
	err := row.Scan(
		&d.ID,
		&d.InvoiceID,
		&d.Type,
		&d.Value,
		&d.ComputedAmount,
		&d.Reason,
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

func (m *DiscountModel) scanRow(rows *sql.Rows) (*Discount, error) {
	d := &Discount{}
	err := rows.Scan(
		&d.ID,
		&d.InvoiceID,
		&d.Type,
		&d.Value,
		&d.ComputedAmount,
		&d.Reason,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	return d, err
}
