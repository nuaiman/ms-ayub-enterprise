package models

import (
	"context"
	"database/sql"
	"time"
)

type Store struct {
	ID              int64      `json:"id"`
	LotID           int64      `json:"lot_id"`
	GodownID        int64      `json:"godown_id"`
	StoreBillType   string     `json:"store_bill_type"` // weight or quantity
	GodownCut       float64    `json:"godown_cut"`
	Quantity        float64    `json:"quantity"`
	QuantityUnit    string     `json:"quantity_unit"`
	Weight          float64    `json:"weight"`
	WeightUnit      string     `json:"weight_unit"`
	IsActive        bool       `json:"is_active"`
	BillingStart    time.Time  `json:"billing_start"`
	BillingEnd      *time.Time `json:"billing_end,omitempty"`
	LastPaidThrough *time.Time `json:"last_paid_through,omitempty"`
	LastPaidAmount  float64    `json:"last_paid_amount"`
	Notes           *string    `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type StoreModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *StoreModel) Insert(ctx context.Context, store *Store) (int64, error) {
	query := `
		INSERT INTO stores (
			lot_id, godown_id, store_bill_type, godown_cut, quantity, quantity_unit,
			weight, weight_unit, is_active, billing_start, billing_end,
			last_paid_through, last_paid_amount, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	isActive := 0
	if store.IsActive {
		isActive = 1
	}

	res, err := m.DB.ExecContext(ctx, query,
		store.LotID,
		store.GodownID,
		store.StoreBillType,
		store.GodownCut,
		store.Quantity,
		store.QuantityUnit,
		store.Weight,
		store.WeightUnit,
		isActive,
		store.BillingStart,
		store.BillingEnd,
		store.LastPaidThrough,
		store.LastPaidAmount,
		store.Notes,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *StoreModel) GetByID(ctx context.Context, id int64) (*Store, error) {
	query := `
		SELECT id, lot_id, godown_id, store_bill_type, godown_cut, quantity, quantity_unit,
		       weight, weight_unit, is_active, billing_start, billing_end,
		       last_paid_through, last_paid_amount, notes, created_at, updated_at
		FROM stores
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanStore(row)
}

func (m *StoreModel) GetByLotID(ctx context.Context, lotID int64) ([]Store, error) {
	query := `
		SELECT id, lot_id, godown_id, store_bill_type, godown_cut, quantity, quantity_unit,
		       weight, weight_unit, is_active, billing_start, billing_end,
		       last_paid_through, last_paid_amount, notes, created_at, updated_at
		FROM stores
		WHERE lot_id = ?
		ORDER BY godown_id ASC, id ASC
	`

	rows, err := m.DB.QueryContext(ctx, query, lotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stores := []Store{}
	for rows.Next() {
		store, err := m.scanStoreRow(rows)
		if err != nil {
			return nil, err
		}
		stores = append(stores, *store)
	}

	return stores, rows.Err()
}

func (m *StoreModel) GetByGodownID(ctx context.Context, godownID int64) ([]Store, error) {
	query := `
		SELECT id, lot_id, godown_id, store_bill_type, godown_cut, quantity, quantity_unit,
		       weight, weight_unit, is_active, billing_start, billing_end,
		       last_paid_through, last_paid_amount, notes, created_at, updated_at
		FROM stores
		WHERE godown_id = ?
		ORDER BY lot_id ASC, id ASC
	`

	rows, err := m.DB.QueryContext(ctx, query, godownID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stores := []Store{}
	for rows.Next() {
		store, err := m.scanStoreRow(rows)
		if err != nil {
			return nil, err
		}
		stores = append(stores, *store)
	}

	return stores, rows.Err()
}

func (m *StoreModel) GetAll(ctx context.Context) ([]Store, error) {
	query := `
		SELECT id, lot_id, godown_id, store_bill_type, godown_cut, quantity, quantity_unit,
		       weight, weight_unit, is_active, billing_start, billing_end,
		       last_paid_through, last_paid_amount, notes, created_at, updated_at
		FROM stores
		ORDER BY lot_id, godown_id ASC, id ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stores := []Store{}
	for rows.Next() {
		store, err := m.scanStoreRow(rows)
		if err != nil {
			return nil, err
		}
		stores = append(stores, *store)
	}

	return stores, rows.Err()
}

func (m *StoreModel) GetActive(ctx context.Context) ([]Store, error) {
	query := `
		SELECT id, lot_id, godown_id, store_bill_type, godown_cut, quantity, quantity_unit,
		       weight, weight_unit, is_active, billing_start, billing_end,
		       last_paid_through, last_paid_amount, notes, created_at, updated_at
		FROM stores
		WHERE is_active = 1
		ORDER BY lot_id, godown_id ASC, id ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stores := []Store{}
	for rows.Next() {
		store, err := m.scanStoreRow(rows)
		if err != nil {
			return nil, err
		}
		stores = append(stores, *store)
	}

	return stores, rows.Err()
}

func (m *StoreModel) GetWithInventory(ctx context.Context) ([]Store, error) {
	query := `
		SELECT id, lot_id, godown_id, store_bill_type, godown_cut, quantity, quantity_unit,
		       weight, weight_unit, is_active, billing_start, billing_end,
		       last_paid_through, last_paid_amount, notes, created_at, updated_at
		FROM stores
		WHERE quantity > 0 OR weight > 0
		ORDER BY lot_id, godown_id ASC, id ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stores := []Store{}
	for rows.Next() {
		store, err := m.scanStoreRow(rows)
		if err != nil {
			return nil, err
		}
		stores = append(stores, *store)
	}

	return stores, rows.Err()
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *StoreModel) Update(ctx context.Context, store *Store) error {
	query := `
		UPDATE stores
		SET 
			store_bill_type = ?,
			godown_cut = ?,
			quantity = ?,
			quantity_unit = ?,
			weight = ?,
			weight_unit = ?,
			is_active = ?,
			billing_start = ?,
			billing_end = ?,
			last_paid_through = ?,
			last_paid_amount = ?,
			notes = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	isActive := 0
	if store.IsActive {
		isActive = 1
	}

	_, err := m.DB.ExecContext(ctx, query,
		store.StoreBillType,
		store.GodownCut,
		store.Quantity,
		store.QuantityUnit,
		store.Weight,
		store.WeightUnit,
		isActive,
		store.BillingStart,
		store.BillingEnd,
		store.LastPaidThrough,
		store.LastPaidAmount,
		store.Notes,
		store.ID,
	)
	return err
}

// UpdateInventory updates a store's quantity/weight. The cascade to
// deactivate the store (and the lot if all its stores are empty) is
// performed by the SQLite triggers on the stores table.
func (m *StoreModel) UpdateInventory(ctx context.Context, id int64, quantity, weight float64) error {
	query := `
		UPDATE stores
		SET 
			quantity = ?,
			weight = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, quantity, weight, id)
	return err
}

func (m *StoreModel) ToggleActive(ctx context.Context, id int64, isActive bool) error {
	active := 0
	if isActive {
		active = 1
	}

	query := `
		UPDATE stores
		SET is_active = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, active, id)
	return err
}

func (m *StoreModel) UpdateBilling(ctx context.Context, id int64, billingStart time.Time, billingEnd *time.Time) error {
	query := `
		UPDATE stores
		SET 
			billing_start = ?,
			billing_end = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, billingStart, billingEnd, id)
	return err
}

func (m *StoreModel) UpdatePayment(ctx context.Context, id int64, lastPaidThrough *time.Time, lastPaidAmount float64) error {
	query := `
		UPDATE stores
		SET 
			last_paid_through = ?,
			last_paid_amount = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, lastPaidThrough, lastPaidAmount, id)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *StoreModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM stores WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *StoreModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM stores WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *StoreModel) GetTotalQuantityByLot(ctx context.Context, lotID int64) (float64, error) {
	query := `SELECT COALESCE(SUM(quantity), 0) FROM stores WHERE lot_id = ?`

	var total float64
	err := m.DB.QueryRowContext(ctx, query, lotID).Scan(&total)
	return total, err
}

func (m *StoreModel) GetTotalWeightByLot(ctx context.Context, lotID int64) (float64, error) {
	query := `SELECT COALESCE(SUM(weight), 0) FROM stores WHERE lot_id = ?`

	var total float64
	err := m.DB.QueryRowContext(ctx, query, lotID).Scan(&total)
	return total, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *StoreModel) scanStore(row *sql.Row) (*Store, error) {
	store := &Store{}
	err := row.Scan(
		&store.ID,
		&store.LotID,
		&store.GodownID,
		&store.StoreBillType,
		&store.GodownCut,
		&store.Quantity,
		&store.QuantityUnit,
		&store.Weight,
		&store.WeightUnit,
		&store.IsActive,
		&store.BillingStart,
		&store.BillingEnd,
		&store.LastPaidThrough,
		&store.LastPaidAmount,
		&store.Notes,
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return store, nil
}

func (m *StoreModel) scanStoreRow(rows *sql.Rows) (*Store, error) {
	store := &Store{}
	err := rows.Scan(
		&store.ID,
		&store.LotID,
		&store.GodownID,
		&store.StoreBillType,
		&store.GodownCut,
		&store.Quantity,
		&store.QuantityUnit,
		&store.Weight,
		&store.WeightUnit,
		&store.IsActive,
		&store.BillingStart,
		&store.BillingEnd,
		&store.LastPaidThrough,
		&store.LastPaidAmount,
		&store.Notes,
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	return store, err
}

func (m *StoreModel) ExistsByLotAndGodown(ctx context.Context, lotID, godownID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM stores WHERE lot_id = ? AND godown_id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, lotID, godownID).Scan(&exists)
	return exists, err
}
