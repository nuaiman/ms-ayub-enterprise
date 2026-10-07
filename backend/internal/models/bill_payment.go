package models

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

type BillPayment struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	BillType        string    `json:"bill_type"` // godown | majhi | customer_store
	BillID          int64     `json:"bill_id"`
	Amount          float64   `json:"amount"`
	PaymentDate     time.Time `json:"payment_date"`
	PaymentMethod   *string   `json:"payment_method,omitempty"`
	ReferenceNumber *string   `json:"reference_number,omitempty"`
	Notes           *string   `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type BillTotals struct {
	TotalPaid        float64
	TotalPaidThrough *time.Time
}

type BillPaymentModel struct {
	DB *sql.DB
}

const billPaymentSelectCols = `
	id, user_id, bill_type, bill_id, amount,
	payment_date, payment_method, reference_number, notes,
	created_at
`

// =============================================================================
// CREATE
// =============================================================================

func (m *BillPaymentModel) Insert(ctx context.Context, p *BillPayment) (int64, error) {
	query := `
		INSERT INTO bill_payments (
			user_id, bill_type, bill_id, amount,
			payment_date, payment_method, reference_number, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		p.UserID,
		p.BillType,
		p.BillID,
		p.Amount,
		p.PaymentDate,
		p.PaymentMethod,
		p.ReferenceNumber,
		p.Notes,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *BillPaymentModel) GetByID(ctx context.Context, id int64) (*BillPayment, error) {
	query := `SELECT ` + billPaymentSelectCols + ` FROM bill_payments WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scan(row)
}

func (m *BillPaymentModel) GetByBill(ctx context.Context, billType string, billID int64) ([]BillPayment, error) {
	query := `SELECT ` + billPaymentSelectCols + `
		FROM bill_payments
		WHERE bill_type = ? AND bill_id = ?
		ORDER BY payment_date DESC, id DESC`
	return m.query(ctx, query, billType, billID)
}

func (m *BillPaymentModel) GetAll(ctx context.Context) ([]BillPayment, error) {
	query := `SELECT ` + billPaymentSelectCols + ` FROM bill_payments ORDER BY id DESC`
	return m.query(ctx, query)
}

// =============================================================================
// AGGREGATES
// =============================================================================

// GetTotalPaid returns the sum of all payments for a bill.
func (m *BillPaymentModel) GetTotalPaid(ctx context.Context, billType string, billID int64) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM bill_payments
		WHERE bill_type = ? AND bill_id = ?
	`
	var total float64
	err := m.DB.QueryRowContext(ctx, query, billType, billID).Scan(&total)
	return total, err
}

// GetTotalPaidThrough returns the timestamp of the most recent payment, or nil if none.
func (m *BillPaymentModel) GetTotalPaidThrough(ctx context.Context, billType string, billID int64) (*time.Time, error) {
	query := `
		SELECT MAX(payment_date)
		FROM bill_payments
		WHERE bill_type = ? AND bill_id = ?
	`
	var ts *time.Time
	err := m.DB.QueryRowContext(ctx, query, billType, billID).Scan(&ts)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return ts, err
}

// GetTotalsForBills returns totals keyed by bill_id, for one bill_type.
// Bills with no payments are omitted from the map (caller treats missing as zero).
func (m *BillPaymentModel) GetTotalsForBills(ctx context.Context, billType string, billIDs []int64) (map[int64]BillTotals, error) {
	if len(billIDs) == 0 {
		return map[int64]BillTotals{}, nil
	}

	placeholders := make([]string, len(billIDs))
	args := make([]any, 0, len(billIDs)+1)
	args = append(args, billType)
	for i, id := range billIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}

	query := `
		SELECT bill_id, COALESCE(SUM(amount), 0), MAX(payment_date)
		FROM bill_payments
		WHERE bill_type = ? AND bill_id IN (` + strings.Join(placeholders, ",") + `)
		GROUP BY bill_id
	`

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[int64]BillTotals, len(billIDs))
	for rows.Next() {
		var (
			billID   int64
			total    float64
			paidThru *time.Time
		)
		if err := rows.Scan(&billID, &total, &paidThru); err != nil {
			return nil, err
		}
		out[billID] = BillTotals{TotalPaid: total, TotalPaidThrough: paidThru}
	}
	return out, rows.Err()
}

// =============================================================================
// DELETE
// =============================================================================

func (m *BillPaymentModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM bill_payments WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *BillPaymentModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM bill_payments WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS / HELPERS
// =============================================================================

func (m *BillPaymentModel) query(ctx context.Context, query string, args ...any) ([]BillPayment, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := []BillPayment{}
	for rows.Next() {
		p, err := m.scanRow(rows)
		if err != nil {
			return nil, err
		}
		payments = append(payments, *p)
	}
	return payments, rows.Err()
}

func (m *BillPaymentModel) scan(row *sql.Row) (*BillPayment, error) {
	p := &BillPayment{}
	err := row.Scan(
		&p.ID,
		&p.UserID,
		&p.BillType,
		&p.BillID,
		&p.Amount,
		&p.PaymentDate,
		&p.PaymentMethod,
		&p.ReferenceNumber,
		&p.Notes,
		&p.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

func (m *BillPaymentModel) scanRow(rows *sql.Rows) (*BillPayment, error) {
	p := &BillPayment{}
	err := rows.Scan(
		&p.ID,
		&p.UserID,
		&p.BillType,
		&p.BillID,
		&p.Amount,
		&p.PaymentDate,
		&p.PaymentMethod,
		&p.ReferenceNumber,
		&p.Notes,
		&p.CreatedAt,
	)
	return p, err
}
