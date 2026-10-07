package models

import (
	"context"
	"database/sql"
	"time"
)

type CustomerAdditionalBill struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CustomerID  int64     `json:"customer_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CustomerAdditionalBillModel struct {
	DB *sql.DB
}

const customerAdditionalBillSelectCols = `
	id, user_id, customer_id, amount, description,
	created_at, updated_at
`

// =============================================================================
// CREATE
// =============================================================================

func (m *CustomerAdditionalBillModel) Insert(ctx context.Context, b *CustomerAdditionalBill) (int64, error) {
	query := `
		INSERT INTO customer_additional_bills (
			user_id, customer_id, amount, description
		) VALUES (?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		b.UserID,
		b.CustomerID,
		b.Amount,
		b.Description,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *CustomerAdditionalBillModel) GetByID(ctx context.Context, id int64) (*CustomerAdditionalBill, error) {
	query := `SELECT ` + customerAdditionalBillSelectCols + ` FROM customer_additional_bills WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scan(row)
}

func (m *CustomerAdditionalBillModel) GetAll(ctx context.Context) ([]CustomerAdditionalBill, error) {
	query := `SELECT ` + customerAdditionalBillSelectCols + ` FROM customer_additional_bills ORDER BY created_at DESC, id DESC`
	return m.query(ctx, query)
}

func (m *CustomerAdditionalBillModel) GetByCustomerID(ctx context.Context, customerID int64) ([]CustomerAdditionalBill, error) {
	query := `SELECT ` + customerAdditionalBillSelectCols + ` FROM customer_additional_bills WHERE customer_id = ? ORDER BY created_at DESC, id DESC`
	return m.query(ctx, query, customerID)
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *CustomerAdditionalBillModel) Update(ctx context.Context, b *CustomerAdditionalBill) error {
	query := `
		UPDATE customer_additional_bills
		SET
			amount = ?,
			description = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := m.DB.ExecContext(ctx, query,
		b.Amount,
		b.Description,
		b.ID,
	)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *CustomerAdditionalBillModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM customer_additional_bills WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *CustomerAdditionalBillModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customer_additional_bills WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS / HELPERS
// =============================================================================

func (m *CustomerAdditionalBillModel) query(ctx context.Context, query string, args ...any) ([]CustomerAdditionalBill, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bills := []CustomerAdditionalBill{}
	for rows.Next() {
		b, err := m.scanRow(rows)
		if err != nil {
			return nil, err
		}
		bills = append(bills, *b)
	}
	return bills, rows.Err()
}

func (m *CustomerAdditionalBillModel) scan(row *sql.Row) (*CustomerAdditionalBill, error) {
	b := &CustomerAdditionalBill{}
	err := row.Scan(
		&b.ID,
		&b.UserID,
		&b.CustomerID,
		&b.Amount,
		&b.Description,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return b, nil
}

func (m *CustomerAdditionalBillModel) scanRow(rows *sql.Rows) (*CustomerAdditionalBill, error) {
	b := &CustomerAdditionalBill{}
	err := rows.Scan(
		&b.ID,
		&b.UserID,
		&b.CustomerID,
		&b.Amount,
		&b.Description,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	return b, err
}
