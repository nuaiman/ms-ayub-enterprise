package models

import (
	"context"
	"database/sql"
	"time"
)

type Lot struct {
	ID                       int64      `json:"id"`
	UserID                   int64      `json:"user_id"`
	CustomerID               *int64     `json:"customer_id,omitempty"`
	ProductName              *string    `json:"product_name,omitempty"`
	Category                 *string    `json:"category,omitempty"`
	LotNumber                int64      `json:"lot_number"`
	CustomerChargeType       string     `json:"customer_charge_type"` // weight or quantity
	MajhiBillType            string     `json:"majhi_bill_type"`      // weight, quantity, or job
	CustomerStorageRate      float64    `json:"customer_storage_rate"`
	UnloadRate               float64    `json:"unload_rate"`
	MajhiID                  *int64     `json:"majhi_id,omitempty"`
	MajhiCut                 float64    `json:"majhi_cut"`
	IsActive                 bool       `json:"is_active"`
	Notes                    *string    `json:"notes,omitempty"`
	ImageURL                 *string    `json:"image_url,omitempty"`
	CustomerLastPaidThrough  *time.Time `json:"customer_last_paid_through,omitempty"`
	CustomerLastPaidAmount   float64    `json:"customer_last_paid_amount"`
	CustomerPaidUnloadAmount float64    `json:"customer_paid_unload_amount"`
	MajhiTotalPaid           float64    `json:"majhi_total_paid"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

type LotModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *LotModel) Insert(ctx context.Context, lot *Lot) (int64, error) {
	query := `
		INSERT INTO lots (
			user_id, customer_id, product_name, category, lot_number,
			customer_charge_type, majhi_bill_type,
			customer_storage_rate, unload_rate, majhi_id, majhi_cut,
			is_active, notes, image_url,
			customer_last_paid_through, customer_last_paid_amount,
			customer_paid_unload_amount, majhi_total_paid
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	isActive := 0
	if lot.IsActive {
		isActive = 1
	}

	res, err := m.DB.ExecContext(ctx, query,
		lot.UserID,
		lot.CustomerID,
		lot.ProductName,
		lot.Category,
		lot.LotNumber,
		lot.CustomerChargeType,
		lot.MajhiBillType,
		lot.CustomerStorageRate,
		lot.UnloadRate,
		lot.MajhiID,
		lot.MajhiCut,
		isActive,
		lot.Notes,
		lot.ImageURL,
		lot.CustomerLastPaidThrough,
		lot.CustomerLastPaidAmount,
		lot.CustomerPaidUnloadAmount,
		lot.MajhiTotalPaid,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *LotModel) GetByID(ctx context.Context, id int64) (*Lot, error) {
	query := `
		SELECT id, user_id, customer_id, product_name, category, lot_number,
		       customer_charge_type, majhi_bill_type,
		       customer_storage_rate, unload_rate, majhi_id, majhi_cut,
		       is_active, notes, image_url,
		       customer_last_paid_through, customer_last_paid_amount,
		       customer_paid_unload_amount, majhi_total_paid,
		       created_at, updated_at
		FROM lots
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanLot(row)
}

func (m *LotModel) GetByCustomerID(ctx context.Context, customerID int64) ([]Lot, error) {
	query := `
		SELECT id, user_id, customer_id, product_name, category, lot_number,
		       customer_charge_type, majhi_bill_type,
		       customer_storage_rate, unload_rate, majhi_id, majhi_cut,
		       is_active, notes, image_url,
		       customer_last_paid_through, customer_last_paid_amount,
		       customer_paid_unload_amount, majhi_total_paid,
		       created_at, updated_at
		FROM lots
		WHERE customer_id = ?
		ORDER BY lot_number ASC
	`

	rows, err := m.DB.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lots := []Lot{}
	for rows.Next() {
		lot, err := m.scanLotRow(rows)
		if err != nil {
			return nil, err
		}
		lots = append(lots, *lot)
	}

	return lots, rows.Err()
}

func (m *LotModel) GetByLotNumber(ctx context.Context, lotNumber int64) (*Lot, error) {
	query := `
		SELECT id, user_id, customer_id, product_name, category, lot_number,
		       customer_charge_type, majhi_bill_type,
		       customer_storage_rate, unload_rate, majhi_id, majhi_cut,
		       is_active, notes, image_url,
		       customer_last_paid_through, customer_last_paid_amount,
		       customer_paid_unload_amount, majhi_total_paid,
		       created_at, updated_at
		FROM lots
		WHERE lot_number = ?
	`

	row := m.DB.QueryRowContext(ctx, query, lotNumber)
	return m.scanLot(row)
}

func (m *LotModel) GetAll(ctx context.Context) ([]Lot, error) {
	query := `
		SELECT id, user_id, customer_id, product_name, category, lot_number,
		       customer_charge_type, majhi_bill_type,
		       customer_storage_rate, unload_rate, majhi_id, majhi_cut,
		       is_active, notes, image_url,
		       customer_last_paid_through, customer_last_paid_amount,
		       customer_paid_unload_amount, majhi_total_paid,
		       created_at, updated_at
		FROM lots
		ORDER BY lot_number ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lots := []Lot{}
	for rows.Next() {
		lot, err := m.scanLotRow(rows)
		if err != nil {
			return nil, err
		}
		lots = append(lots, *lot)
	}

	return lots, rows.Err()
}

func (m *LotModel) GetActive(ctx context.Context) ([]Lot, error) {
	query := `
		SELECT id, user_id, customer_id, product_name, category, lot_number,
		       customer_charge_type, majhi_bill_type,
		       customer_storage_rate, unload_rate, majhi_id, majhi_cut,
		       is_active, notes, image_url,
		       customer_last_paid_through, customer_last_paid_amount,
		       customer_paid_unload_amount, majhi_total_paid,
		       created_at, updated_at
		FROM lots
		WHERE is_active = 1
		ORDER BY lot_number ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lots := []Lot{}
	for rows.Next() {
		lot, err := m.scanLotRow(rows)
		if err != nil {
			return nil, err
		}
		lots = append(lots, *lot)
	}

	return lots, rows.Err()
}

func (m *LotModel) Search(ctx context.Context, query string) ([]Lot, error) {
	searchQuery := `
		SELECT id, user_id, customer_id, product_name, category, lot_number,
		       customer_charge_type, majhi_bill_type,
		       customer_storage_rate, unload_rate, majhi_id, majhi_cut,
		       is_active, notes, image_url,
		       customer_last_paid_through, customer_last_paid_amount,
		       customer_paid_unload_amount, majhi_total_paid,
		       created_at, updated_at
		FROM lots
		WHERE product_name LIKE ? OR category LIKE ? OR lot_number LIKE ?
		ORDER BY lot_number ASC
	`

	searchTerm := "%" + query + "%"

	rows, err := m.DB.QueryContext(ctx, searchQuery, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lots := []Lot{}
	for rows.Next() {
		lot, err := m.scanLotRow(rows)
		if err != nil {
			return nil, err
		}
		lots = append(lots, *lot)
	}

	return lots, rows.Err()
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *LotModel) Update(ctx context.Context, lot *Lot) error {
	query := `
		UPDATE lots
		SET 
			customer_id = ?,
			product_name = ?,
			category = ?,
			customer_charge_type = ?,
			majhi_bill_type = ?,
			customer_storage_rate = ?,
			unload_rate = ?,
			majhi_id = ?,
			majhi_cut = ?,
			is_active = ?,
			notes = ?,
			customer_last_paid_through = ?,
			customer_last_paid_amount = ?,
			customer_paid_unload_amount = ?,
			majhi_total_paid = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	isActive := 0
	if lot.IsActive {
		isActive = 1
	}

	_, err := m.DB.ExecContext(ctx, query,
		lot.CustomerID,
		lot.ProductName,
		lot.Category,
		lot.CustomerChargeType,
		lot.MajhiBillType,
		lot.CustomerStorageRate,
		lot.UnloadRate,
		lot.MajhiID,
		lot.MajhiCut,
		isActive,
		lot.Notes,
		lot.CustomerLastPaidThrough,
		lot.CustomerLastPaidAmount,
		lot.CustomerPaidUnloadAmount,
		lot.MajhiTotalPaid,
		lot.ID,
	)
	return err
}

func (m *LotModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	query := `
		UPDATE lots
		SET image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, imageURL, id)
	return err
}

func (m *LotModel) ToggleActive(ctx context.Context, id int64, isActive bool) error {
	active := 0
	if isActive {
		active = 1
	}

	query := `
		UPDATE lots
		SET is_active = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, active, id)
	return err
}

// UpdateCustomerPayment updates customer payment tracking fields
func (m *LotModel) UpdateCustomerPayment(ctx context.Context, id int64, paidThrough *time.Time, paidAmount float64) error {
	query := `
		UPDATE lots
		SET 
			customer_last_paid_through = ?,
			customer_last_paid_amount = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, paidThrough, paidAmount, id)
	return err
}

// UpdateCustomerUnloadPayment updates customer unload payment amount
func (m *LotModel) UpdateCustomerUnloadPayment(ctx context.Context, id int64, paidUnloadAmount float64) error {
	query := `
		UPDATE lots
		SET 
			customer_paid_unload_amount = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, paidUnloadAmount, id)
	return err
}

// UpdateMajhiPayment updates majhi payment tracking
func (m *LotModel) UpdateMajhiPayment(ctx context.Context, id int64, totalPaid float64) error {
	query := `
		UPDATE lots
		SET 
			majhi_total_paid = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, totalPaid, id)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *LotModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM lots WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *LotModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM lots WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// GetTotalUnloadAmountByCustomer returns total unload bill amount for a customer
func (m *LotModel) GetTotalUnloadAmountByCustomer(ctx context.Context, customerID int64) (float64, error) {
	query := `
		SELECT COALESCE(SUM(l.unload_rate * 
			CASE 
				WHEN l.customer_charge_type = 'quantity' THEN COALESCE(s.quantity, 0)
				ELSE COALESCE(s.weight, 0)
			END
		), 0)
		FROM lots l
		INNER JOIN stores s ON s.lot_id = l.id
		WHERE l.customer_id = ? AND l.is_active = 1 AND l.unload_rate > 0
	`

	var total float64
	err := m.DB.QueryRowContext(ctx, query, customerID).Scan(&total)
	return total, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *LotModel) scanLot(row *sql.Row) (*Lot, error) {
	lot := &Lot{}
	var isActive int
	err := row.Scan(
		&lot.ID,
		&lot.UserID,
		&lot.CustomerID,
		&lot.ProductName,
		&lot.Category,
		&lot.LotNumber,
		&lot.CustomerChargeType,
		&lot.MajhiBillType,
		&lot.CustomerStorageRate,
		&lot.UnloadRate,
		&lot.MajhiID,
		&lot.MajhiCut,
		&isActive,
		&lot.Notes,
		&lot.ImageURL,
		&lot.CustomerLastPaidThrough,
		&lot.CustomerLastPaidAmount,
		&lot.CustomerPaidUnloadAmount,
		&lot.MajhiTotalPaid,
		&lot.CreatedAt,
		&lot.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	lot.IsActive = isActive == 1
	return lot, nil
}

func (m *LotModel) scanLotRow(rows *sql.Rows) (*Lot, error) {
	lot := &Lot{}
	var isActive int
	err := rows.Scan(
		&lot.ID,
		&lot.UserID,
		&lot.CustomerID,
		&lot.ProductName,
		&lot.Category,
		&lot.LotNumber,
		&lot.CustomerChargeType,
		&lot.MajhiBillType,
		&lot.CustomerStorageRate,
		&lot.UnloadRate,
		&lot.MajhiID,
		&lot.MajhiCut,
		&isActive,
		&lot.Notes,
		&lot.ImageURL,
		&lot.CustomerLastPaidThrough,
		&lot.CustomerLastPaidAmount,
		&lot.CustomerPaidUnloadAmount,
		&lot.MajhiTotalPaid,
		&lot.CreatedAt,
		&lot.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	lot.IsActive = isActive == 1
	return lot, nil
}

func (m *LotModel) ExistsByCustomerAndLot(ctx context.Context, customerID, lotNumber int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM lots WHERE customer_id = ? AND lot_number = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, customerID, lotNumber).Scan(&exists)
	return exists, err
}
