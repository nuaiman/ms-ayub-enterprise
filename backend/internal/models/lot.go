package models

import (
	"context"
	"database/sql"
	"time"
)

type Lot struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	CustomerID   int64     `json:"customer_id"`
	LotNumber    string    `json:"lot_number"`
	ProductName  string    `json:"product_name"`
	WeightUnit   string    `json:"weight_unit"`
	QuantityUnit string    `json:"quantity_unit"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
			user_id, customer_id, lot_number, product_name,
			weight_unit, quantity_unit
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		lot.UserID,
		lot.CustomerID,
		lot.LotNumber,
		lot.ProductName,
		lot.WeightUnit,
		lot.QuantityUnit,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

const lotSelectCols = `
	id, user_id, customer_id, lot_number, product_name,
	weight_unit, quantity_unit,
	created_at, updated_at
`

func (m *LotModel) GetByID(ctx context.Context, id int64) (*Lot, error) {
	query := `SELECT ` + lotSelectCols + ` FROM lots WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanLot(row)
}

func (m *LotModel) GetByCustomerID(ctx context.Context, customerID int64) ([]Lot, error) {
	query := `SELECT ` + lotSelectCols + ` FROM lots WHERE customer_id = ? ORDER BY lot_number ASC`
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

func (m *LotModel) GetAll(ctx context.Context) ([]Lot, error) {
	query := `SELECT ` + lotSelectCols + ` FROM lots ORDER BY lot_number ASC`
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
		SELECT ` + lotSelectCols + `
		FROM lots
		WHERE product_name LIKE ? OR lot_number LIKE ?
		ORDER BY lot_number ASC
	`
	searchTerm := "%" + query + "%"

	rows, err := m.DB.QueryContext(ctx, searchQuery, searchTerm, searchTerm)
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
			lot_number = ?,
			product_name = ?,
			weight_unit = ?,
			quantity_unit = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		lot.CustomerID,
		lot.LotNumber,
		lot.ProductName,
		lot.WeightUnit,
		lot.QuantityUnit,
		lot.ID,
	)
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

func (m *LotModel) ExistsByCustomerAndLot(ctx context.Context, customerID int64, lotNumber string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM lots WHERE customer_id = ? AND lot_number = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, customerID, lotNumber).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *LotModel) scanLot(row *sql.Row) (*Lot, error) {
	lot := &Lot{}
	err := row.Scan(
		&lot.ID,
		&lot.UserID,
		&lot.CustomerID,
		&lot.LotNumber,
		&lot.ProductName,
		&lot.WeightUnit,
		&lot.QuantityUnit,
		&lot.CreatedAt,
		&lot.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return lot, nil
}

func (m *LotModel) scanLotRow(rows *sql.Rows) (*Lot, error) {
	lot := &Lot{}
	err := rows.Scan(
		&lot.ID,
		&lot.UserID,
		&lot.CustomerID,
		&lot.LotNumber,
		&lot.ProductName,
		&lot.WeightUnit,
		&lot.QuantityUnit,
		&lot.CreatedAt,
		&lot.UpdatedAt,
	)
	return lot, err
}
