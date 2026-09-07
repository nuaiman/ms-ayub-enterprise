package models

import (
	"context"
	"database/sql"
	"time"
)

type Lot struct {
	ID                       int64      `json:"id"`
	ItemID                   int64      `json:"item_id"`
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
			item_id, lot_number, customer_charge_type, majhi_bill_type,
			customer_storage_rate, unload_rate, majhi_id, majhi_cut,
			is_active, notes, image_url,
			customer_last_paid_through, customer_last_paid_amount, customer_paid_unload_amount, majhi_total_paid
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	isActive := 0
	if lot.IsActive {
		isActive = 1
	}

	res, err := m.DB.ExecContext(ctx, query,
		lot.ItemID,
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
		SELECT id, item_id, lot_number, customer_charge_type, majhi_bill_type,
		       customer_storage_rate, unload_rate, majhi_id, majhi_cut,
		       is_active, notes, image_url,
		       customer_last_paid_through, customer_last_paid_amount, customer_paid_unload_amount, majhi_total_paid,
		       created_at, updated_at
		FROM lots
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanLot(row)
}

func (m *LotModel) GetByItemID(ctx context.Context, itemID int64) ([]Lot, error) {
	query := `
		SELECT id, item_id, lot_number, customer_charge_type, majhi_bill_type,
		       customer_storage_rate, unload_rate, majhi_id, majhi_cut,
		       is_active, notes, image_url,
		       customer_last_paid_through, customer_last_paid_amount, customer_paid_unload_amount, majhi_total_paid,
		       created_at, updated_at
		FROM lots
		WHERE item_id = ?
		ORDER BY lot_number ASC
	`

	rows, err := m.DB.QueryContext(ctx, query, itemID)
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

func (m *LotModel) GetByItemAndLotNumber(ctx context.Context, itemID int64, lotNumber int64) (*Lot, error) {
	query := `
		SELECT id, item_id, lot_number, customer_charge_type, majhi_bill_type,
		       customer_storage_rate, unload_rate, majhi_id, majhi_cut,
		       is_active, notes, image_url,
		       customer_last_paid_through, customer_last_paid_amount, customer_paid_unload_amount, majhi_total_paid,
		       created_at, updated_at
		FROM lots
		WHERE item_id = ? AND lot_number = ?
	`

	row := m.DB.QueryRowContext(ctx, query, itemID, lotNumber)
	return m.scanLot(row)
}

func (m *LotModel) GetAll(ctx context.Context) ([]Lot, error) {
	query := `
		SELECT id, item_id, lot_number, customer_charge_type, majhi_bill_type,
		       customer_storage_rate, unload_rate, majhi_id, majhi_cut,
		       is_active, notes, image_url,
		       customer_last_paid_through, customer_last_paid_amount, customer_paid_unload_amount, majhi_total_paid,
		       created_at, updated_at
		FROM lots
		ORDER BY item_id, lot_number ASC
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
		SELECT id, item_id, lot_number, customer_charge_type, majhi_bill_type,
		       customer_storage_rate, unload_rate, majhi_id, majhi_cut,
		       is_active, notes, image_url,
		       customer_last_paid_through, customer_last_paid_amount, customer_paid_unload_amount, majhi_total_paid,
		       created_at, updated_at
		FROM lots
		WHERE is_active = 1
		ORDER BY item_id, lot_number ASC
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
		SELECT l.id, l.item_id, l.lot_number, l.customer_charge_type, l.majhi_bill_type,
		       l.customer_storage_rate, l.unload_rate, l.majhi_id, l.majhi_cut,
		       l.is_active, l.notes, l.image_url,
		       l.customer_last_paid_through, l.customer_last_paid_amount, l.customer_paid_unload_amount, l.majhi_total_paid,
		       l.created_at, l.updated_at
		FROM lots l
		INNER JOIN items i ON l.item_id = i.id
		WHERE i.product_name LIKE ? OR i.category LIKE ? OR l.lot_number LIKE ?
		ORDER BY l.item_id, l.lot_number ASC
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

func (m *LotModel) ExistsByItemAndLot(ctx context.Context, itemID int64, lotNumber int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM lots WHERE item_id = ? AND lot_number = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, itemID, lotNumber).Scan(&exists)
	return exists, err
}

func (m *LotModel) CountByItem(ctx context.Context, itemID int64) (int, error) {
	query := `SELECT COUNT(*) FROM lots WHERE item_id = ?`

	var count int
	err := m.DB.QueryRowContext(ctx, query, itemID).Scan(&count)
	return count, err
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
		INNER JOIN items i ON l.item_id = i.id
		INNER JOIN stores s ON s.lot_id = l.id
		WHERE i.customer_id = ? AND l.is_active = 1 AND l.unload_rate > 0
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
		&lot.ItemID,
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
		&lot.ItemID,
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
