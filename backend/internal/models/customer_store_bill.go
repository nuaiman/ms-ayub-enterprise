package models

import (
	"context"
	"database/sql"
	"time"
)

type CustomerStoreBill struct {
	ID                    int64     `json:"id"`
	UserID                int64     `json:"user_id"`
	CustomerID            int64     `json:"customer_id"`
	StoreID               int64     `json:"store_id"`
	MonthYear             string    `json:"month_year"` // YYYY-MM
	BillType              string    `json:"bill_type"`  // weight | quantity
	Rate                  float64   `json:"rate"`
	WeightAtBilling       float64   `json:"weight_at_billing"`
	QuantityAtBilling     float64   `json:"quantity_at_billing"`
	WeightUnitAtBilling   *string   `json:"weight_unit_at_billing,omitempty"`
	QuantityUnitAtBilling *string   `json:"quantity_unit_at_billing,omitempty"`
	TotalAmount           float64   `json:"total_amount"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type CustomerStoreBillModel struct {
	DB *sql.DB
}

const customerStoreBillSelectCols = `
	id, user_id, customer_id, store_id, month_year,
	bill_type, rate,
	weight_at_billing, quantity_at_billing,
	weight_unit_at_billing, quantity_unit_at_billing,
	total_amount,
	created_at, updated_at
`

// =============================================================================
// CREATE
// =============================================================================

func (m *CustomerStoreBillModel) Insert(ctx context.Context, bill *CustomerStoreBill) (int64, error) {
	query := `
		INSERT INTO customer_store_bills (
			user_id, customer_id, store_id, month_year,
			bill_type, rate,
			weight_at_billing, quantity_at_billing,
			weight_unit_at_billing, quantity_unit_at_billing,
			total_amount
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		bill.UserID,
		bill.CustomerID,
		bill.StoreID,
		bill.MonthYear,
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

func (m *CustomerStoreBillModel) GetByID(ctx context.Context, id int64) (*CustomerStoreBill, error) {
	query := `SELECT ` + customerStoreBillSelectCols + ` FROM customer_store_bills WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanCustomerStoreBill(row)
}

func (m *CustomerStoreBillModel) GetAll(ctx context.Context) ([]CustomerStoreBill, error) {
	query := `SELECT ` + customerStoreBillSelectCols + ` FROM customer_store_bills ORDER BY month_year DESC, created_at DESC`
	return m.queryCustomerStoreBills(ctx, query)
}

func (m *CustomerStoreBillModel) GetByStoreID(ctx context.Context, storeID int64) ([]CustomerStoreBill, error) {
	query := `SELECT ` + customerStoreBillSelectCols + ` FROM customer_store_bills WHERE store_id = ? ORDER BY month_year DESC`
	return m.queryCustomerStoreBills(ctx, query, storeID)
}

func (m *CustomerStoreBillModel) GetByCustomerID(ctx context.Context, customerID int64) ([]CustomerStoreBill, error) {
	query := `SELECT ` + customerStoreBillSelectCols + ` FROM customer_store_bills WHERE customer_id = ? ORDER BY month_year DESC`
	return m.queryCustomerStoreBills(ctx, query, customerID)
}

func (m *CustomerStoreBillModel) GetByMonth(ctx context.Context, monthYear string) ([]CustomerStoreBill, error) {
	query := `SELECT ` + customerStoreBillSelectCols + ` FROM customer_store_bills WHERE month_year = ? ORDER BY customer_id ASC, store_id ASC`
	return m.queryCustomerStoreBills(ctx, query, monthYear)
}

func (m *CustomerStoreBillModel) GetByCustomerAndMonth(ctx context.Context, customerID int64, monthYear string) ([]CustomerStoreBill, error) {
	query := `SELECT ` + customerStoreBillSelectCols + ` FROM customer_store_bills WHERE customer_id = ? AND month_year = ? ORDER BY store_id ASC`
	return m.queryCustomerStoreBills(ctx, query, customerID, monthYear)
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *CustomerStoreBillModel) Update(ctx context.Context, bill *CustomerStoreBill) error {
	query := `
		UPDATE customer_store_bills
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

func (m *CustomerStoreBillModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM customer_store_bills WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *CustomerStoreBillModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customer_store_bills WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *CustomerStoreBillModel) ExistsForCustomerStoreMonth(
	ctx context.Context,
	customerID, storeID int64,
	monthYear string,
) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customer_store_bills WHERE customer_id = ? AND store_id = ? AND month_year = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, customerID, storeID, monthYear).Scan(&exists)
	return exists, err
}

// EnsureCurrentMonthBills creates a draft customer store bill for every
// (active store × its customer) that doesn't already have one for the
// current month. Rate starts at 0 — set per bill via Update afterwards.
func (m *CustomerStoreBillModel) EnsureCurrentMonthBills(ctx context.Context, userID int64) error {
	currentMonth := time.Now().Format("2006-01")

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Join active stores to their lot -> customer so we know who to bill.
	rows, err := tx.QueryContext(ctx, `
		SELECT s.id, l.customer_id, s.weight, s.quantity, l.weight_unit, l.quantity_unit
		FROM stores s
		JOIN lots l ON l.id = s.lot_id
		WHERE s.is_active = 1
	`)
	if err != nil {
		return err
	}

	type pair struct {
		storeID      int64
		customerID   int64
		weight       float64
		quantity     float64
		weightUnit   string
		quantityUnit string
	}
	pairs := []pair{}
	for rows.Next() {
		var p pair
		if err := rows.Scan(&p.storeID, &p.customerID, &p.weight, &p.quantity, &p.weightUnit, &p.quantityUnit); err != nil {
			rows.Close()
			return err
		}
		pairs = append(pairs, p)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT OR IGNORE INTO customer_store_bills (
			user_id, customer_id, store_id, month_year,
			bill_type, rate,
			weight_at_billing, quantity_at_billing,
			weight_unit_at_billing, quantity_unit_at_billing,
			total_amount
		) VALUES (?, ?, ?, ?, 'quantity', 0, ?, ?, ?, ?, 0)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range pairs {
		if _, err := stmt.ExecContext(ctx,
			userID, p.customerID, p.storeID, currentMonth,
			p.weight, p.quantity,
			p.weightUnit, p.quantityUnit,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// =============================================================================
// SCANNERS / HELPERS
// =============================================================================

func (m *CustomerStoreBillModel) queryCustomerStoreBills(ctx context.Context, query string, args ...any) ([]CustomerStoreBill, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bills := []CustomerStoreBill{}
	for rows.Next() {
		b, err := m.scanCustomerStoreBillRow(rows)
		if err != nil {
			return nil, err
		}
		bills = append(bills, *b)
	}
	return bills, rows.Err()
}

func (m *CustomerStoreBillModel) scanCustomerStoreBill(row *sql.Row) (*CustomerStoreBill, error) {
	b := &CustomerStoreBill{}
	err := row.Scan(
		&b.ID,
		&b.UserID,
		&b.CustomerID,
		&b.StoreID,
		&b.MonthYear,
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

func (m *CustomerStoreBillModel) scanCustomerStoreBillRow(rows *sql.Rows) (*CustomerStoreBill, error) {
	b := &CustomerStoreBill{}
	err := rows.Scan(
		&b.ID,
		&b.UserID,
		&b.CustomerID,
		&b.StoreID,
		&b.MonthYear,
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
