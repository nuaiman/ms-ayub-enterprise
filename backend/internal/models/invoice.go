package models

import (
	"context"
	"database/sql"
	"time"
)

type Invoice struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	EntityType     string    `json:"entity_type"`
	EntityID       int64     `json:"entity_id"`
	Subtotal       float64   `json:"subtotal"`
	DiscountAmount float64   `json:"discount_amount"`
	Total          float64   `json:"total"`
	Notes          *string   `json:"notes,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type InvoiceModel struct {
	DB *sql.DB
}

const invoiceSelectCols = `
	id, user_id, entity_type, entity_id,
	subtotal, discount_amount, total,
	notes,
	created_at, updated_at
`

// =============================================================================
// CREATE
// =============================================================================

func (m *InvoiceModel) Insert(ctx context.Context, inv *Invoice) (int64, error) {
	query := `
		INSERT INTO invoices (
			user_id, entity_type, entity_id,
			subtotal, discount_amount, total,
			notes
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		inv.UserID,
		inv.EntityType,
		inv.EntityID,
		inv.Subtotal,
		inv.DiscountAmount,
		inv.Total,
		inv.Notes,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *InvoiceModel) GetByID(ctx context.Context, id int64) (*Invoice, error) {
	query := `SELECT ` + invoiceSelectCols + ` FROM invoices WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scan(row)
}

func (m *InvoiceModel) GetAll(ctx context.Context) ([]Invoice, error) {
	query := `SELECT ` + invoiceSelectCols + ` FROM invoices ORDER BY created_at DESC, id DESC`
	return m.query(ctx, query)
}

func (m *InvoiceModel) GetByEntity(ctx context.Context, entityType string, entityID int64) ([]Invoice, error) {
	query := `SELECT ` + invoiceSelectCols + `
		FROM invoices
		WHERE entity_type = ? AND entity_id = ?
		ORDER BY created_at DESC, id DESC`
	return m.query(ctx, query, entityType, entityID)
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *InvoiceModel) Update(ctx context.Context, inv *Invoice) error {
	query := `
		UPDATE invoices
		SET
			subtotal = ?,
			discount_amount = ?,
			total = ?,
			notes = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := m.DB.ExecContext(ctx, query,
		inv.Subtotal,
		inv.DiscountAmount,
		inv.Total,
		inv.Notes,
		inv.ID,
	)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *InvoiceModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM invoices WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *InvoiceModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM invoices WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS / HELPERS
// =============================================================================

func (m *InvoiceModel) query(ctx context.Context, query string, args ...any) ([]Invoice, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := []Invoice{}
	for rows.Next() {
		inv, err := m.scanRow(rows)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, *inv)
	}
	return invoices, rows.Err()
}

func (m *InvoiceModel) scan(row *sql.Row) (*Invoice, error) {
	inv := &Invoice{}
	err := row.Scan(
		&inv.ID,
		&inv.UserID,
		&inv.EntityType,
		&inv.EntityID,
		&inv.Subtotal,
		&inv.DiscountAmount,
		&inv.Total,
		&inv.Notes,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return inv, nil
}

func (m *InvoiceModel) scanRow(rows *sql.Rows) (*Invoice, error) {
	inv := &Invoice{}
	err := rows.Scan(
		&inv.ID,
		&inv.UserID,
		&inv.EntityType,
		&inv.EntityID,
		&inv.Subtotal,
		&inv.DiscountAmount,
		&inv.Total,
		&inv.Notes,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	)
	return inv, err
}
