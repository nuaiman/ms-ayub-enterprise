package models

import (
	"context"
	"database/sql"
	"time"
)

type AdditionalCharge struct {
	ID                       int64      `json:"id"`
	UserID                   int64      `json:"user_id"`
	CustomerID               int64      `json:"customer_id"`
	Amount                   float64    `json:"amount"`
	Description              string     `json:"description"`
	CustomerTotalPaid        float64    `json:"customer_total_paid"`
	CustomerTotalPaidThrough *time.Time `json:"customer_total_paid_through,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

type AdditionalChargeModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *AdditionalChargeModel) Insert(ctx context.Context, charge *AdditionalCharge) (int64, error) {
	query := `
		INSERT INTO customer_additional_charges (
			user_id, customer_id,
			amount, description,
			customer_total_paid, customer_total_paid_through
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		charge.UserID,
		charge.CustomerID,
		charge.Amount,
		charge.Description,
		charge.CustomerTotalPaid,
		charge.CustomerTotalPaidThrough,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

const additionalChargeSelectCols = `
	id, user_id, customer_id,
	amount, description,
	customer_total_paid, customer_total_paid_through,
	created_at, updated_at
`

func (m *AdditionalChargeModel) GetByID(ctx context.Context, id int64) (*AdditionalCharge, error) {
	query := `SELECT ` + additionalChargeSelectCols + ` FROM customer_additional_charges WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanAdditionalCharge(row)
}

func (m *AdditionalChargeModel) GetAll(ctx context.Context) ([]AdditionalCharge, error) {
	query := `SELECT ` + additionalChargeSelectCols + ` FROM customer_additional_charges ORDER BY created_at DESC`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	charges := []AdditionalCharge{}
	for rows.Next() {
		charge, err := m.scanAdditionalChargeRow(rows)
		if err != nil {
			return nil, err
		}
		charges = append(charges, *charge)
	}

	return charges, rows.Err()
}

func (m *AdditionalChargeModel) GetByCustomerID(ctx context.Context, customerID int64) ([]AdditionalCharge, error) {
	query := `SELECT ` + additionalChargeSelectCols + ` FROM customer_additional_charges WHERE customer_id = ? ORDER BY created_at DESC`
	rows, err := m.DB.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	charges := []AdditionalCharge{}
	for rows.Next() {
		charge, err := m.scanAdditionalChargeRow(rows)
		if err != nil {
			return nil, err
		}
		charges = append(charges, *charge)
	}

	return charges, rows.Err()
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *AdditionalChargeModel) Update(ctx context.Context, charge *AdditionalCharge) error {
	query := `
		UPDATE customer_additional_charges
		SET
			amount = ?,
			description = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		charge.Amount,
		charge.Description,
		charge.ID,
	)
	return err
}

func (m *AdditionalChargeModel) UpdateCustomerPayment(ctx context.Context, id int64, totalPaid float64, paidThrough *time.Time) error {
	query := `
		UPDATE customer_additional_charges
		SET
			customer_total_paid = ?,
			customer_total_paid_through = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, totalPaid, paidThrough, id)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *AdditionalChargeModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM customer_additional_charges WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *AdditionalChargeModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customer_additional_charges WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *AdditionalChargeModel) scanAdditionalCharge(row *sql.Row) (*AdditionalCharge, error) {
	charge := &AdditionalCharge{}
	err := row.Scan(
		&charge.ID,
		&charge.UserID,
		&charge.CustomerID,
		&charge.Amount,
		&charge.Description,
		&charge.CustomerTotalPaid,
		&charge.CustomerTotalPaidThrough,
		&charge.CreatedAt,
		&charge.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return charge, nil
}

func (m *AdditionalChargeModel) scanAdditionalChargeRow(rows *sql.Rows) (*AdditionalCharge, error) {
	charge := &AdditionalCharge{}
	err := rows.Scan(
		&charge.ID,
		&charge.UserID,
		&charge.CustomerID,
		&charge.Amount,
		&charge.Description,
		&charge.CustomerTotalPaid,
		&charge.CustomerTotalPaidThrough,
		&charge.CreatedAt,
		&charge.UpdatedAt,
	)
	return charge, err
}
