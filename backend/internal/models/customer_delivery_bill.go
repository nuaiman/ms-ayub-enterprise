package models

import (
	"context"
	"database/sql"
	"time"
)

type CustomerDeliveryBill struct {
	ID                    int64     `json:"id"`
	UserID                int64     `json:"user_id"`
	CustomerID            int64     `json:"customer_id"`
	DeliveryItemID        int64     `json:"delivery_item_id"`
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

type CustomerDeliveryBillModel struct {
	DB *sql.DB
}

const customerDeliveryBillSelectCols = `
	id, user_id, customer_id, delivery_item_id,
	bill_type, rate,
	weight_at_billing, quantity_at_billing,
	weight_unit_at_billing, quantity_unit_at_billing,
	total_amount,
	created_at, updated_at
`

func (m *CustomerDeliveryBillModel) Insert(ctx context.Context, bill *CustomerDeliveryBill) (int64, error) {
	query := `
		INSERT INTO customer_delivery_bills (
			user_id, customer_id, delivery_item_id,
			bill_type, rate,
			weight_at_billing, quantity_at_billing,
			weight_unit_at_billing, quantity_unit_at_billing,
			total_amount
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		bill.UserID,
		bill.CustomerID,
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

func (m *CustomerDeliveryBillModel) GetByID(ctx context.Context, id int64) (*CustomerDeliveryBill, error) {
	query := `SELECT ` + customerDeliveryBillSelectCols + ` FROM customer_delivery_bills WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scan(row)
}

func (m *CustomerDeliveryBillModel) GetAll(ctx context.Context) ([]CustomerDeliveryBill, error) {
	query := `SELECT ` + customerDeliveryBillSelectCols + ` FROM customer_delivery_bills ORDER BY created_at DESC`
	return m.query(ctx, query)
}

func (m *CustomerDeliveryBillModel) GetByDeliveryItemID(ctx context.Context, deliveryItemID int64) (*CustomerDeliveryBill, error) {
	query := `SELECT ` + customerDeliveryBillSelectCols + ` FROM customer_delivery_bills WHERE delivery_item_id = ?`
	row := m.DB.QueryRowContext(ctx, query, deliveryItemID)
	return m.scan(row)
}

func (m *CustomerDeliveryBillModel) GetByCustomerID(ctx context.Context, customerID int64) ([]CustomerDeliveryBill, error) {
	query := `SELECT ` + customerDeliveryBillSelectCols + ` FROM customer_delivery_bills WHERE customer_id = ? ORDER BY created_at DESC`
	return m.query(ctx, query, customerID)
}

func (m *CustomerDeliveryBillModel) Update(ctx context.Context, bill *CustomerDeliveryBill) error {
	query := `
		UPDATE customer_delivery_bills
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

func (m *CustomerDeliveryBillModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM customer_delivery_bills WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m *CustomerDeliveryBillModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customer_delivery_bills WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *CustomerDeliveryBillModel) ExistsForDeliveryItem(ctx context.Context, deliveryItemID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customer_delivery_bills WHERE delivery_item_id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, deliveryItemID).Scan(&exists)
	return exists, err
}

func (m *CustomerDeliveryBillModel) query(ctx context.Context, query string, args ...any) ([]CustomerDeliveryBill, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bills := []CustomerDeliveryBill{}
	for rows.Next() {
		b, err := m.scanRow(rows)
		if err != nil {
			return nil, err
		}
		bills = append(bills, *b)
	}
	return bills, rows.Err()
}

func (m *CustomerDeliveryBillModel) scan(row *sql.Row) (*CustomerDeliveryBill, error) {
	b := &CustomerDeliveryBill{}
	err := row.Scan(
		&b.ID,
		&b.UserID,
		&b.CustomerID,
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

func (m *CustomerDeliveryBillModel) scanRow(rows *sql.Rows) (*CustomerDeliveryBill, error) {
	b := &CustomerDeliveryBill{}
	err := rows.Scan(
		&b.ID,
		&b.UserID,
		&b.CustomerID,
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
