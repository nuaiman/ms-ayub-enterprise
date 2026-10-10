package models

import (
	"context"
	"database/sql"
	"time"
)

type CustomerUnloadBill struct {
	ID                    int64     `json:"id"`
	UserID                int64     `json:"user_id"`
	CustomerID            int64     `json:"customer_id"`
	StoreID               int64     `json:"store_id"`
	BillType              string    `json:"bill_type"`
	Rate                  float64   `json:"rate"`
	WeightAtBilling       float64   `json:"weight_at_billing"`
	QuantityAtBilling     float64   `json:"quantity_at_billing"`
	WeightUnitAtBilling   *string   `json:"weight_unit_at_billing,omitempty"`
	QuantityUnitAtBilling *string   `json:"quantity_unit_at_billing,omitempty"`
	TotalAmount           float64   `json:"total_amount"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type CustomerUnloadBillModel struct {
	DB *sql.DB
}

const customerUnloadBillSelectCols = `
	id, user_id, customer_id, store_id,
	bill_type, rate,
	weight_at_billing, quantity_at_billing,
	weight_unit_at_billing, quantity_unit_at_billing,
	total_amount,
	created_at, updated_at
`

// =============================================================================
// CREATE
// =============================================================================

func (m *CustomerUnloadBillModel) Insert(ctx context.Context, bill *CustomerUnloadBill) (int64, error) {
	query := `
		INSERT INTO customer_unload_bills (
			user_id, customer_id, store_id,
			bill_type, rate,
			weight_at_billing, quantity_at_billing,
			weight_unit_at_billing, quantity_unit_at_billing,
			total_amount
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		bill.UserID,
		bill.CustomerID,
		bill.StoreID,
		bill.BillType,
		bill.Rate,
		bill.WeightAtBilling,
		bill.QuantityAtBilling,
		bill.WeightUnitAtBilling,
		bill.QuantityUnitAtBilling,
		bill.TotalAmount,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *CustomerUnloadBillModel) GetByID(ctx context.Context, id int64) (*CustomerUnloadBill, error) {
	query := `SELECT ` + customerUnloadBillSelectCols + ` FROM customer_unload_bills WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scan(row)
}

func (m *CustomerUnloadBillModel) GetAll(ctx context.Context) ([]CustomerUnloadBill, error) {
	query := `SELECT ` + customerUnloadBillSelectCols + ` FROM customer_unload_bills ORDER BY created_at DESC, id DESC`
	return m.query(ctx, query)
}

func (m *CustomerUnloadBillModel) GetByCustomerID(ctx context.Context, customerID int64) ([]CustomerUnloadBill, error) {
	query := `SELECT ` + customerUnloadBillSelectCols + ` FROM customer_unload_bills WHERE customer_id = ? ORDER BY created_at DESC, id DESC`
	return m.query(ctx, query, customerID)
}

func (m *CustomerUnloadBillModel) GetByStoreID(ctx context.Context, storeID int64) ([]CustomerUnloadBill, error) {
	query := `SELECT ` + customerUnloadBillSelectCols + ` FROM customer_unload_bills WHERE store_id = ? ORDER BY created_at DESC, id DESC`
	return m.query(ctx, query, storeID)
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *CustomerUnloadBillModel) Update(ctx context.Context, bill *CustomerUnloadBill) error {
	query := `
		UPDATE customer_unload_bills
		SET
			bill_type = ?,
			rate = ?,
			weight_at_billing = ?,
			quantity_at_billing = ?,
			total_amount = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := m.DB.ExecContext(ctx, query,
		bill.BillType,
		bill.Rate,
		bill.WeightAtBilling,
		bill.QuantityAtBilling,
		bill.TotalAmount,
		bill.ID,
	)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *CustomerUnloadBillModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM customer_unload_bills WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *CustomerUnloadBillModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customer_unload_bills WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS / HELPERS
// =============================================================================

func (m *CustomerUnloadBillModel) query(ctx context.Context, query string, args ...any) ([]CustomerUnloadBill, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bills := []CustomerUnloadBill{}
	for rows.Next() {
		b, err := m.scanRow(rows)
		if err != nil {
			return nil, err
		}
		bills = append(bills, *b)
	}
	return bills, rows.Err()
}

func (m *CustomerUnloadBillModel) scan(row *sql.Row) (*CustomerUnloadBill, error) {
	b := &CustomerUnloadBill{}
	err := row.Scan(
		&b.ID,
		&b.UserID,
		&b.CustomerID,
		&b.StoreID,
		&b.BillType,
		&b.Rate,
		&b.WeightAtBilling,
		&b.QuantityAtBilling,
		&b.WeightUnitAtBilling,
		&b.QuantityUnitAtBilling,
		&b.TotalAmount,
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

func (m *CustomerUnloadBillModel) scanRow(rows *sql.Rows) (*CustomerUnloadBill, error) {
	b := &CustomerUnloadBill{}
	err := rows.Scan(
		&b.ID,
		&b.UserID,
		&b.CustomerID,
		&b.StoreID,
		&b.BillType,
		&b.Rate,
		&b.WeightAtBilling,
		&b.QuantityAtBilling,
		&b.WeightUnitAtBilling,
		&b.QuantityUnitAtBilling,
		&b.TotalAmount,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	return b, err
}