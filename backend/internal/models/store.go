package models

import (
	"context"
	"database/sql"
	"time"
)

type Store struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	LotID     int64     `json:"lot_id"`
	GodownID  int64     `json:"godown_id"`
	Weight    float64   `json:"weight"`
	Quantity  float64   `json:"quantity"`
	StartDate time.Time `json:"start_date"`
	IsActive  bool      `json:"is_active"`
	ImageURL  *string   `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
			user_id, lot_id, godown_id, weight, quantity, start_date,
			is_active,
			image_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		store.UserID,
		store.LotID,
		store.GodownID,
		store.Weight,
		store.Quantity,
		store.StartDate,
		boolToInt(store.IsActive),
		store.ImageURL,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

const storeSelectCols = `
	id, user_id, lot_id, godown_id, weight, quantity, start_date,
	is_active,
	image_url,
	created_at, updated_at
`

func (m *StoreModel) GetByID(ctx context.Context, id int64) (*Store, error) {
	query := `SELECT ` + storeSelectCols + ` FROM stores WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanStore(row)
}

func (m *StoreModel) GetByLotID(ctx context.Context, lotID int64) ([]Store, error) {
	query := `SELECT ` + storeSelectCols + ` FROM stores WHERE lot_id = ? ORDER BY id ASC`
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

func (m *StoreModel) GetAll(ctx context.Context) ([]Store, error) {
	query := `SELECT ` + storeSelectCols + ` FROM stores ORDER BY lot_id ASC, id ASC`
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
	query := `SELECT ` + storeSelectCols + ` FROM stores WHERE is_active = 1 ORDER BY lot_id ASC, id ASC`
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
	query := `SELECT ` + storeSelectCols + ` FROM stores WHERE quantity > 0 OR weight > 0 ORDER BY lot_id ASC, id ASC`
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
			weight = ?,
			quantity = ?,
			start_date = ?,
			is_active = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		store.Weight,
		store.Quantity,
		store.StartDate,
		boolToInt(store.IsActive),
		store.ID,
	)
	return err
}

// UpdateInventory updates a store's quantity/weight. Cascade deactivation
// is handled by SQLite triggers.
func (m *StoreModel) UpdateInventory(ctx context.Context, id int64, quantity, weight float64) error {
	query := `
		UPDATE stores
		SET quantity = ?, weight = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := m.DB.ExecContext(ctx, query, quantity, weight, id)
	return err
}

func (m *StoreModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	query := `
		UPDATE stores
		SET image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := m.DB.ExecContext(ctx, query, imageURL, id)
	return err
}

func (m *StoreModel) ToggleActive(ctx context.Context, id int64, isActive bool) error {
	query := `UPDATE stores SET is_active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, boolToInt(isActive), id)
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

func (m *StoreModel) ExistsByLot(ctx context.Context, lotID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM stores WHERE lot_id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, lotID).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *StoreModel) scanStore(row *sql.Row) (*Store, error) {
	store := &Store{}
	var isActive int
	err := row.Scan(
		&store.ID,
		&store.UserID,
		&store.LotID,
		&store.GodownID,
		&store.Weight,
		&store.Quantity,
		&store.StartDate,
		&isActive,
		&store.ImageURL,
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	store.IsActive = isActive == 1
	return store, nil
}

func (m *StoreModel) scanStoreRow(rows *sql.Rows) (*Store, error) {
	store := &Store{}
	var isActive int
	err := rows.Scan(
		&store.ID,
		&store.UserID,
		&store.LotID,
		&store.GodownID,
		&store.Weight,
		&store.Quantity,
		&store.StartDate,
		&isActive,
		&store.ImageURL,
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	store.IsActive = isActive == 1
	return store, nil
}

// UpdateGodown moves a store to a different godown. Used by transfers.
func (m *StoreModel) UpdateGodown(ctx context.Context, id int64, godownID int64) error {
	query := `
		UPDATE stores
		SET godown_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := m.DB.ExecContext(ctx, query, godownID, id)
	return err
}

// =============================================================================
// HELPERS
// =============================================================================

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
