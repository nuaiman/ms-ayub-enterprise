package models

import (
	"context"
	"database/sql"
	"time"
)

type MajhiBill struct {
	ID                    int64     `json:"id"`
	UserID                int64     `json:"user_id"`
	MajhiID               int64     `json:"majhi_id"`
	StoreID               *int64    `json:"store_id,omitempty"`
	DeliveryItemID        *int64    `json:"delivery_item_id,omitempty"`
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

type MajhiBillModel struct {
	DB *sql.DB
}

const majhiBillSelectCols = `
	id, user_id, majhi_id,
	store_id, delivery_item_id,
	bill_type, rate,
	weight_at_billing, quantity_at_billing,
	weight_unit_at_billing, quantity_unit_at_billing,
	total_amount,
	created_at, updated_at
`

func (m *MajhiBillModel) Insert(ctx context.Context, bill *MajhiBill) (int64, error) {
	query := `
		INSERT INTO majhi_bills (
			user_id, majhi_id,
			store_id, delivery_item_id,
			bill_type, rate,
			weight_at_billing, quantity_at_billing,
			weight_unit_at_billing, quantity_unit_at_billing,
			total_amount
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		bill.UserID,
		bill.MajhiID,
		bill.StoreID,
		bill.DeliveryItemID,
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

func (m *MajhiBillModel) GetByID(ctx context.Context, id int64) (*MajhiBill, error) {
	query := `SELECT ` + majhiBillSelectCols + ` FROM majhi_bills WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanMajhiBill(row)
}

func (m *MajhiBillModel) GetAll(ctx context.Context) ([]MajhiBill, error) {
	query := `SELECT ` + majhiBillSelectCols + ` FROM majhi_bills ORDER BY created_at DESC`
	return m.queryMajhiBills(ctx, query)
}

func (m *MajhiBillModel) GetByStoreID(ctx context.Context, storeID int64) ([]MajhiBill, error) {
	query := `SELECT ` + majhiBillSelectCols + ` FROM majhi_bills WHERE store_id = ? ORDER BY created_at DESC`
	return m.queryMajhiBills(ctx, query, storeID)
}

func (m *MajhiBillModel) GetByDeliveryItemID(ctx context.Context, deliveryItemID int64) ([]MajhiBill, error) {
	query := `SELECT ` + majhiBillSelectCols + ` FROM majhi_bills WHERE delivery_item_id = ? ORDER BY created_at DESC`
	return m.queryMajhiBills(ctx, query, deliveryItemID)
}

func (m *MajhiBillModel) GetByMajhiID(ctx context.Context, majhiID int64) ([]MajhiBill, error) {
	query := `SELECT ` + majhiBillSelectCols + ` FROM majhi_bills WHERE majhi_id = ? ORDER BY created_at DESC`
	return m.queryMajhiBills(ctx, query, majhiID)
}

func (m *MajhiBillModel) Update(ctx context.Context, bill *MajhiBill) error {
	query := `
		UPDATE majhi_bills
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

func (m *MajhiBillModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM majhi_bills WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m *MajhiBillModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM majhi_bills WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *MajhiBillModel) ExistsForStore(ctx context.Context, storeID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM majhi_bills WHERE store_id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, storeID).Scan(&exists)
	return exists, err
}

func (m *MajhiBillModel) ExistsForDeliveryItem(ctx context.Context, deliveryItemID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM majhi_bills WHERE delivery_item_id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, deliveryItemID).Scan(&exists)
	return exists, err
}

func (m *MajhiBillModel) queryMajhiBills(ctx context.Context, query string, args ...any) ([]MajhiBill, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bills := []MajhiBill{}
	for rows.Next() {
		b, err := m.scanMajhiBillRow(rows)
		if err != nil {
			return nil, err
		}
		bills = append(bills, *b)
	}
	return bills, rows.Err()
}

func (m *MajhiBillModel) scanMajhiBill(row *sql.Row) (*MajhiBill, error) {
	b := &MajhiBill{}
	err := row.Scan(
		&b.ID,
		&b.UserID,
		&b.MajhiID,
		&b.StoreID,
		&b.DeliveryItemID,
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

func (m *MajhiBillModel) scanMajhiBillRow(rows *sql.Rows) (*MajhiBill, error) {
	b := &MajhiBill{}
	err := rows.Scan(
		&b.ID,
		&b.UserID,
		&b.MajhiID,
		&b.StoreID,
		&b.DeliveryItemID,
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
