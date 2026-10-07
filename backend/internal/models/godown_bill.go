package models

import (
	"context"
	"database/sql"
	"time"
)

type GodownBill struct {
	ID                    int64     `json:"id"`
	UserID                int64     `json:"user_id"`
	GodownID              int64     `json:"godown_id"`
	StoreID               int64     `json:"store_id"`
	MonthYear             string    `json:"month_year"` // YYYY-MM
	BillType              string    `json:"bill_type"`  // weight | quantity | fixed
	Rate                  float64   `json:"rate"`
	WeightAtBilling       float64   `json:"weight_at_billing"`
	QuantityAtBilling     float64   `json:"quantity_at_billing"`
	WeightUnitAtBilling   *string   `json:"weight_unit_at_billing,omitempty"`
	QuantityUnitAtBilling *string   `json:"quantity_unit_at_billing,omitempty"`
	TotalAmount           float64   `json:"total_amount"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type GodownBillModel struct {
	DB *sql.DB
}

const godownBillSelectCols = `
	id, user_id, godown_id, store_id, month_year,
	bill_type, rate,
	weight_at_billing, quantity_at_billing,
	weight_unit_at_billing, quantity_unit_at_billing,
	total_amount,
	created_at, updated_at
`

// =============================================================================
// CREATE
// =============================================================================

func (m *GodownBillModel) Insert(ctx context.Context, bill *GodownBill) (int64, error) {
	query := `
		INSERT INTO godown_bills (
			user_id, godown_id, store_id, month_year,
			bill_type, rate,
			weight_at_billing, quantity_at_billing,
			weight_unit_at_billing, quantity_unit_at_billing,
			total_amount
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		bill.UserID,
		bill.GodownID,
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

func (m *GodownBillModel) GetByID(ctx context.Context, id int64) (*GodownBill, error) {
	query := `SELECT ` + godownBillSelectCols + ` FROM godown_bills WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanGodownBill(row)
}

func (m *GodownBillModel) GetAll(ctx context.Context) ([]GodownBill, error) {
	query := `SELECT ` + godownBillSelectCols + ` FROM godown_bills ORDER BY month_year DESC, created_at DESC`
	return m.queryGodownBills(ctx, query)
}

func (m *GodownBillModel) GetByStoreID(ctx context.Context, storeID int64) ([]GodownBill, error) {
	query := `SELECT ` + godownBillSelectCols + ` FROM godown_bills WHERE store_id = ? ORDER BY month_year DESC`
	return m.queryGodownBills(ctx, query, storeID)
}

func (m *GodownBillModel) GetByGodownID(ctx context.Context, godownID int64) ([]GodownBill, error) {
	query := `SELECT ` + godownBillSelectCols + ` FROM godown_bills WHERE godown_id = ? ORDER BY month_year DESC`
	return m.queryGodownBills(ctx, query, godownID)
}

func (m *GodownBillModel) GetByMonth(ctx context.Context, monthYear string) ([]GodownBill, error) {
	query := `SELECT ` + godownBillSelectCols + ` FROM godown_bills WHERE month_year = ? ORDER BY godown_id ASC, store_id ASC`
	return m.queryGodownBills(ctx, query, monthYear)
}

func (m *GodownBillModel) GetByGodownAndMonth(ctx context.Context, godownID int64, monthYear string) ([]GodownBill, error) {
	query := `SELECT ` + godownBillSelectCols + ` FROM godown_bills WHERE godown_id = ? AND month_year = ? ORDER BY store_id ASC`
	return m.queryGodownBills(ctx, query, godownID, monthYear)
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *GodownBillModel) Update(ctx context.Context, bill *GodownBill) error {
	query := `
		UPDATE godown_bills
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

func (m *GodownBillModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM godown_bills WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *GodownBillModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM godown_bills WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *GodownBillModel) ExistsForGodownStoreMonth(ctx context.Context, godownID, storeID int64, monthYear string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM godown_bills WHERE godown_id = ? AND store_id = ? AND month_year = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, godownID, storeID, monthYear).Scan(&exists)
	return exists, err
}

// EnsureCurrentMonthBills creates a draft godown bill for every
// (active store × active godown) that doesn't already have one for the
// current month. Rate starts at 0 — set per bill via Update afterwards.
func (m *GodownBillModel) EnsureCurrentMonthBills(ctx context.Context, userID int64) error {
	currentMonth := time.Now().Format("2006-01")

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Fetch active stores with their current weight/quantity + lot units
	rows, err := tx.QueryContext(ctx, `
		SELECT s.id, s.weight, s.quantity, l.weight_unit, l.quantity_unit
		FROM stores s
		JOIN lots l ON l.id = s.lot_id
		WHERE s.is_active = 1
	`)
	if err != nil {
		return err
	}

	type storeInfo struct {
		id           int64
		weight       float64
		quantity     float64
		weightUnit   string
		quantityUnit string
	}
	stores := []storeInfo{}
	for rows.Next() {
		var s storeInfo
		if err := rows.Scan(&s.id, &s.weight, &s.quantity, &s.weightUnit, &s.quantityUnit); err != nil {
			rows.Close()
			return err
		}
		stores = append(stores, s)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	// Fetch active godowns
	grows, err := tx.QueryContext(ctx, `SELECT id FROM godowns`)
	if err != nil {
		return err
	}
	godownIDs := []int64{}
	for grows.Next() {
		var id int64
		if err := grows.Scan(&id); err != nil {
			grows.Close()
			return err
		}
		godownIDs = append(godownIDs, id)
	}
	grows.Close()
	if err := grows.Err(); err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT OR IGNORE INTO godown_bills (
			user_id, godown_id, store_id, month_year,
			bill_type, rate,
			weight_at_billing, quantity_at_billing,
			weight_unit_at_billing, quantity_unit_at_billing,
			total_amount
		) VALUES (?, ?, ?, ?, 'fixed', 0, ?, ?, ?, ?, 0)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, gid := range godownIDs {
		for _, s := range stores {
			if _, err := stmt.ExecContext(ctx,
				userID, gid, s.id, currentMonth,
				s.weight, s.quantity,
				s.weightUnit, s.quantityUnit,
			); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// =============================================================================
// SCANNERS / HELPERS
// =============================================================================

func (m *GodownBillModel) queryGodownBills(ctx context.Context, query string, args ...any) ([]GodownBill, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bills := []GodownBill{}
	for rows.Next() {
		b, err := m.scanGodownBillRow(rows)
		if err != nil {
			return nil, err
		}
		bills = append(bills, *b)
	}
	return bills, rows.Err()
}

func (m *GodownBillModel) scanGodownBill(row *sql.Row) (*GodownBill, error) {
	b := &GodownBill{}
	err := row.Scan(
		&b.ID,
		&b.UserID,
		&b.GodownID,
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

func (m *GodownBillModel) scanGodownBillRow(rows *sql.Rows) (*GodownBill, error) {
	b := &GodownBill{}
	err := rows.Scan(
		&b.ID,
		&b.UserID,
		&b.GodownID,
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
