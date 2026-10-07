package models

import (
	"context"
	"database/sql"
	"time"
)

type InvoiceItem struct {
	ID        int64     `json:"id"`
	InvoiceID int64     `json:"invoice_id"`
	BillType  string    `json:"bill_type"`
	BillID    int64     `json:"bill_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type InvoiceItemModel struct {
	DB *sql.DB
}

const invoiceItemSelectCols = `
	id, invoice_id, bill_type, bill_id, amount,
	created_at
`

// =============================================================================
// CREATE
// =============================================================================

func (m *InvoiceItemModel) Insert(ctx context.Context, item *InvoiceItem) (int64, error) {
	query := `
		INSERT INTO invoice_items (
			invoice_id, bill_type, bill_id, amount
		) VALUES (?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		item.InvoiceID,
		item.BillType,
		item.BillID,
		item.Amount,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// InsertMany inserts all items in a single transaction.
func (m *InvoiceItemModel) InsertMany(ctx context.Context, items []*InvoiceItem) error {
	if len(items) == 0 {
		return nil
	}

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO invoice_items (
			invoice_id, bill_type, bill_id, amount
		) VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		if _, err := stmt.ExecContext(ctx,
			item.InvoiceID,
			item.BillType,
			item.BillID,
			item.Amount,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// =============================================================================
// READ
// =============================================================================

func (m *InvoiceItemModel) GetByID(ctx context.Context, id int64) (*InvoiceItem, error) {
	query := `SELECT ` + invoiceItemSelectCols + ` FROM invoice_items WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scan(row)
}

func (m *InvoiceItemModel) GetByInvoiceID(ctx context.Context, invoiceID int64) ([]InvoiceItem, error) {
	query := `SELECT ` + invoiceItemSelectCols + `
		FROM invoice_items
		WHERE invoice_id = ?
		ORDER BY id ASC`
	return m.query(ctx, query, invoiceID)
}

// ExistsForBill returns true if any invoice_items row references this bill.
// Used to block deletion of a bill that's already been invoiced.
func (m *InvoiceItemModel) ExistsForBill(ctx context.Context, billType string, billID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM invoice_items WHERE bill_type = ? AND bill_id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, billType, billID).Scan(&exists)
	return exists, err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *InvoiceItemModel) DeleteByInvoiceID(ctx context.Context, invoiceID int64) error {
	query := `DELETE FROM invoice_items WHERE invoice_id = ?`
	_, err := m.DB.ExecContext(ctx, query, invoiceID)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *InvoiceItemModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM invoice_items WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS / HELPERS
// =============================================================================

func (m *InvoiceItemModel) query(ctx context.Context, query string, args ...any) ([]InvoiceItem, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []InvoiceItem{}
	for rows.Next() {
		it, err := m.scanRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *it)
	}
	return items, rows.Err()
}

func (m *InvoiceItemModel) scan(row *sql.Row) (*InvoiceItem, error) {
	it := &InvoiceItem{}
	err := row.Scan(
		&it.ID,
		&it.InvoiceID,
		&it.BillType,
		&it.BillID,
		&it.Amount,
		&it.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return it, nil
}

func (m *InvoiceItemModel) scanRow(rows *sql.Rows) (*InvoiceItem, error) {
	it := &InvoiceItem{}
	err := rows.Scan(
		&it.ID,
		&it.InvoiceID,
		&it.BillType,
		&it.BillID,
		&it.Amount,
		&it.CreatedAt,
	)
	return it, err
}
