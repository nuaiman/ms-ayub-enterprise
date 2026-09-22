package models

import (
	"context"
	"database/sql"
	"time"
)

type DeliveryItem struct {
	ID                       int64     `json:"id"`
	DeliveryID               int64     `json:"delivery_id"`
	StoreID                  int64     `json:"store_id"`
	MajhiID                  *int64    `json:"majhi_id,omitempty"`
	LotID                    int64     `json:"lot_id"`
	VehicleNumber            *string   `json:"vehicle_number,omitempty"`
	DriverNumber             *string   `json:"driver_number,omitempty"`
	Quantity                 float64   `json:"quantity"`
	QuantityUnit             string    `json:"quantity_unit"`
	Weight                   float64   `json:"weight"`
	WeightUnit               string    `json:"weight_unit"`
	LoadingRate              float64   `json:"loading_rate"`
	MajhiCut                 float64   `json:"majhi_cut"`
	Notes                    *string   `json:"notes,omitempty"`
	CustomerChargeType       string    `json:"customer_charge_type"`
	CustomerPaidUnloadAmount float64   `json:"customer_paid_unload_amount"`
	MajhiBillType            string    `json:"majhi_bill_type"`
	MajhiTotalPaid           float64   `json:"majhi_total_paid"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type DeliveryItemModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *DeliveryItemModel) Insert(ctx context.Context, item *DeliveryItem) (int64, error) {
	query := `
		INSERT INTO delivery_items (
			delivery_id, store_id, majhi_id, lot_id,
			vehicle_number, driver_number, quantity, quantity_unit,
			weight, weight_unit, loading_rate, majhi_cut, notes,
			customer_charge_type, customer_paid_unload_amount,
			majhi_bill_type, majhi_total_paid
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		item.DeliveryID,
		item.StoreID,
		item.MajhiID,
		item.LotID,
		item.VehicleNumber,
		item.DriverNumber,
		item.Quantity,
		item.QuantityUnit,
		item.Weight,
		item.WeightUnit,
		item.LoadingRate,
		item.MajhiCut,
		item.Notes,
		item.CustomerChargeType,
		item.CustomerPaidUnloadAmount,
		item.MajhiBillType,
		item.MajhiTotalPaid,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *DeliveryItemModel) GetByID(ctx context.Context, id int64) (*DeliveryItem, error) {
	query := `
		SELECT id, delivery_id, store_id, majhi_id, lot_id,
		       vehicle_number, driver_number, quantity, quantity_unit,
		       weight, weight_unit, loading_rate, majhi_cut, notes,
		       customer_charge_type, customer_paid_unload_amount,
		       majhi_bill_type, majhi_total_paid,
		       created_at, updated_at
		FROM delivery_items
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanDeliveryItem(row)
}

func (m *DeliveryItemModel) GetByDeliveryID(ctx context.Context, deliveryID int64) ([]DeliveryItem, error) {
	query := `
		SELECT id, delivery_id, store_id, majhi_id, lot_id,
		       vehicle_number, driver_number, quantity, quantity_unit,
		       weight, weight_unit, loading_rate, majhi_cut, notes,
		       customer_charge_type, customer_paid_unload_amount,
		       majhi_bill_type, majhi_total_paid,
		       created_at, updated_at
		FROM delivery_items
		WHERE delivery_id = ?
		ORDER BY id ASC
	`

	rows, err := m.DB.QueryContext(ctx, query, deliveryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []DeliveryItem{}
	for rows.Next() {
		item, err := m.scanDeliveryItemRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, rows.Err()
}

func (m *DeliveryItemModel) GetByStoreID(ctx context.Context, storeID int64) ([]DeliveryItem, error) {
	query := `
		SELECT id, delivery_id, store_id, majhi_id, lot_id,
		       vehicle_number, driver_number, quantity, quantity_unit,
		       weight, weight_unit, loading_rate, majhi_cut, notes,
		       customer_charge_type, customer_paid_unload_amount,
		       majhi_bill_type, majhi_total_paid,
		       created_at, updated_at
		FROM delivery_items
		WHERE store_id = ?
		ORDER BY created_at DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []DeliveryItem{}
	for rows.Next() {
		item, err := m.scanDeliveryItemRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, rows.Err()
}

func (m *DeliveryItemModel) GetByLotID(ctx context.Context, lotID int64) ([]DeliveryItem, error) {
	query := `
		SELECT id, delivery_id, store_id, majhi_id, lot_id,
		       vehicle_number, driver_number, quantity, quantity_unit,
		       weight, weight_unit, loading_rate, majhi_cut, notes,
		       customer_charge_type, customer_paid_unload_amount,
		       majhi_bill_type, majhi_total_paid,
		       created_at, updated_at
		FROM delivery_items
		WHERE lot_id = ?
		ORDER BY created_at DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, lotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []DeliveryItem{}
	for rows.Next() {
		item, err := m.scanDeliveryItemRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, rows.Err()
}

func (m *DeliveryItemModel) GetAll(ctx context.Context) ([]DeliveryItem, error) {
	query := `
		SELECT id, delivery_id, store_id, majhi_id, lot_id,
		       vehicle_number, driver_number, quantity, quantity_unit,
		       weight, weight_unit, loading_rate, majhi_cut, notes,
		       customer_charge_type, customer_paid_unload_amount,
		       majhi_bill_type, majhi_total_paid,
		       created_at, updated_at
		FROM delivery_items
		ORDER BY created_at DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []DeliveryItem{}
	for rows.Next() {
		item, err := m.scanDeliveryItemRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, rows.Err()
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *DeliveryItemModel) Update(ctx context.Context, item *DeliveryItem) error {
	query := `
		UPDATE delivery_items
		SET 
			store_id = ?,
			majhi_id = ?,
			vehicle_number = ?,
			driver_number = ?,
			quantity = ?,
			quantity_unit = ?,
			weight = ?,
			weight_unit = ?,
			loading_rate = ?,
			majhi_cut = ?,
			notes = ?,
			customer_charge_type = ?,
			customer_paid_unload_amount = ?,
			majhi_bill_type = ?,
			majhi_total_paid = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		item.StoreID,
		item.MajhiID,
		item.VehicleNumber,
		item.DriverNumber,
		item.Quantity,
		item.QuantityUnit,
		item.Weight,
		item.WeightUnit,
		item.LoadingRate,
		item.MajhiCut,
		item.Notes,
		item.CustomerChargeType,
		item.CustomerPaidUnloadAmount,
		item.MajhiBillType,
		item.MajhiTotalPaid,
		item.ID,
	)
	return err
}

// UpdateCustomerUnloadPayment updates the customer unload payment amount for a delivery item
func (m *DeliveryItemModel) UpdateCustomerUnloadPayment(ctx context.Context, id int64, paidAmount float64) error {
	query := `
		UPDATE delivery_items
		SET 
			customer_paid_unload_amount = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, paidAmount, id)
	return err
}

// UpdateMajhiPayment updates the majhi payment amount for a delivery item
func (m *DeliveryItemModel) UpdateMajhiPayment(ctx context.Context, id int64, paidAmount float64) error {
	query := `
		UPDATE delivery_items
		SET 
			majhi_total_paid = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, paidAmount, id)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *DeliveryItemModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM delivery_items WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m *DeliveryItemModel) DeleteByDeliveryID(ctx context.Context, deliveryID int64) error {
	query := `DELETE FROM delivery_items WHERE delivery_id = ?`
	_, err := m.DB.ExecContext(ctx, query, deliveryID)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *DeliveryItemModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM delivery_items WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *DeliveryItemModel) GetTotalQuantityByDelivery(ctx context.Context, deliveryID int64) (float64, error) {
	query := `SELECT COALESCE(SUM(quantity), 0) FROM delivery_items WHERE delivery_id = ?`

	var total float64
	err := m.DB.QueryRowContext(ctx, query, deliveryID).Scan(&total)
	return total, err
}

func (m *DeliveryItemModel) GetTotalWeightByDelivery(ctx context.Context, deliveryID int64) (float64, error) {
	query := `SELECT COALESCE(SUM(weight), 0) FROM delivery_items WHERE delivery_id = ?`

	var total float64
	err := m.DB.QueryRowContext(ctx, query, deliveryID).Scan(&total)
	return total, err
}

// GetCustomerUnloadBillTotal calculates total unload bill for a customer from delivery items
func (m *DeliveryItemModel) GetCustomerUnloadBillTotal(ctx context.Context, customerID int64) (float64, error) {
	query := `
		SELECT COALESCE(SUM(
			CASE 
				WHEN di.customer_charge_type = 'quantity' THEN di.quantity * di.loading_rate
				ELSE di.weight * di.loading_rate
			END
		), 0)
		FROM delivery_items di
		INNER JOIN deliveries d ON di.delivery_id = d.id
		WHERE d.customer_id = ?
	`

	var total float64
	err := m.DB.QueryRowContext(ctx, query, customerID).Scan(&total)
	return total, err
}

// GetMajhiBillTotal calculates total majhi bill for a majhi from delivery items
func (m *DeliveryItemModel) GetMajhiBillTotal(ctx context.Context, majhiID int64) (float64, error) {
	query := `
		SELECT COALESCE(SUM(
			CASE 
				WHEN di.majhi_bill_type = 'quantity' THEN di.quantity * di.majhi_cut
				WHEN di.majhi_bill_type = 'weight' THEN di.weight * di.majhi_cut
				ELSE di.majhi_cut
			END
		), 0)
		FROM delivery_items di
		WHERE di.majhi_id = ?
	`

	var total float64
	err := m.DB.QueryRowContext(ctx, query, majhiID).Scan(&total)
	return total, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *DeliveryItemModel) scanDeliveryItem(row *sql.Row) (*DeliveryItem, error) {
	item := &DeliveryItem{}
	err := row.Scan(
		&item.ID,
		&item.DeliveryID,
		&item.StoreID,
		&item.MajhiID,
		&item.LotID,
		&item.VehicleNumber,
		&item.DriverNumber,
		&item.Quantity,
		&item.QuantityUnit,
		&item.Weight,
		&item.WeightUnit,
		&item.LoadingRate,
		&item.MajhiCut,
		&item.Notes,
		&item.CustomerChargeType,
		&item.CustomerPaidUnloadAmount,
		&item.MajhiBillType,
		&item.MajhiTotalPaid,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (m *DeliveryItemModel) scanDeliveryItemRow(rows *sql.Rows) (*DeliveryItem, error) {
	item := &DeliveryItem{}
	err := rows.Scan(
		&item.ID,
		&item.DeliveryID,
		&item.StoreID,
		&item.MajhiID,
		&item.LotID,
		&item.VehicleNumber,
		&item.DriverNumber,
		&item.Quantity,
		&item.QuantityUnit,
		&item.Weight,
		&item.WeightUnit,
		&item.LoadingRate,
		&item.MajhiCut,
		&item.Notes,
		&item.CustomerChargeType,
		&item.CustomerPaidUnloadAmount,
		&item.MajhiBillType,
		&item.MajhiTotalPaid,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	return item, err
}
