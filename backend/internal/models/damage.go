package models

import (
	"context"
	"database/sql"
	"time"
)

type Damage struct {
	ID           int64     `json:"id"`
	StoreID      int64     `json:"store_id"`
	UserID       int64     `json:"user_id"`
	Quantity     float64   `json:"quantity"`
	QuantityUnit string    `json:"quantity_unit"`
	Weight       float64   `json:"weight"`
	WeightUnit   string    `json:"weight_unit"`
	DamageDate   time.Time `json:"damage_date"`
	Reason       string    `json:"reason"`
	Amount       float64   `json:"amount"`
	Notes        *string   `json:"notes,omitempty"`
	ImageURL     *string   `json:"image_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type DamageModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *DamageModel) Insert(ctx context.Context, damage *Damage) (int64, error) {
	query := `
		INSERT INTO damages (
			store_id, user_id, quantity, quantity_unit, weight, weight_unit,
			damage_date, reason, amount, notes, image_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		damage.StoreID,
		damage.UserID,
		damage.Quantity,
		damage.QuantityUnit,
		damage.Weight,
		damage.WeightUnit,
		damage.DamageDate,
		damage.Reason,
		damage.Amount,
		damage.Notes,
		damage.ImageURL,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *DamageModel) GetByID(ctx context.Context, id int64) (*Damage, error) {
	query := `
		SELECT id, store_id, user_id, quantity, quantity_unit, weight, weight_unit,
		       damage_date, reason, amount, notes, image_url, created_at, updated_at
		FROM damages
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanDamage(row)
}

func (m *DamageModel) GetByStoreID(ctx context.Context, storeID int64) ([]Damage, error) {
	query := `
		SELECT id, store_id, user_id, quantity, quantity_unit, weight, weight_unit,
		       damage_date, reason, amount, notes, image_url, created_at, updated_at
		FROM damages
		WHERE store_id = ?
		ORDER BY damage_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	damages := []Damage{}
	for rows.Next() {
		damage, err := m.scanDamageRow(rows)
		if err != nil {
			return nil, err
		}
		damages = append(damages, *damage)
	}

	return damages, rows.Err()
}

func (m *DamageModel) GetByUserID(ctx context.Context, userID int64) ([]Damage, error) {
	query := `
		SELECT id, store_id, user_id, quantity, quantity_unit, weight, weight_unit,
		       damage_date, reason, amount, notes, image_url, created_at, updated_at
		FROM damages
		WHERE user_id = ?
		ORDER BY damage_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	damages := []Damage{}
	for rows.Next() {
		damage, err := m.scanDamageRow(rows)
		if err != nil {
			return nil, err
		}
		damages = append(damages, *damage)
	}

	return damages, rows.Err()
}

func (m *DamageModel) GetAll(ctx context.Context) ([]Damage, error) {
	query := `
		SELECT id, store_id, user_id, quantity, quantity_unit, weight, weight_unit,
		       damage_date, reason, amount, notes, image_url, created_at, updated_at
		FROM damages
		ORDER BY damage_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	damages := []Damage{}
	for rows.Next() {
		damage, err := m.scanDamageRow(rows)
		if err != nil {
			return nil, err
		}
		damages = append(damages, *damage)
	}

	return damages, rows.Err()
}

func (m *DamageModel) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]Damage, error) {
	query := `
		SELECT id, store_id, user_id, quantity, quantity_unit, weight, weight_unit,
		       damage_date, reason, amount, notes, image_url, created_at, updated_at
		FROM damages
		WHERE damage_date >= ? AND damage_date <= ?
		ORDER BY damage_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	damages := []Damage{}
	for rows.Next() {
		damage, err := m.scanDamageRow(rows)
		if err != nil {
			return nil, err
		}
		damages = append(damages, *damage)
	}

	return damages, rows.Err()
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *DamageModel) Update(ctx context.Context, damage *Damage) error {
	query := `
		UPDATE damages
		SET 
			quantity = ?,
			quantity_unit = ?,
			weight = ?,
			weight_unit = ?,
			damage_date = ?,
			reason = ?,
			amount = ?,
			notes = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		damage.Quantity,
		damage.QuantityUnit,
		damage.Weight,
		damage.WeightUnit,
		damage.DamageDate,
		damage.Reason,
		damage.Amount,
		damage.Notes,
		damage.ID,
	)
	return err
}

func (m *DamageModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	query := `
		UPDATE damages
		SET image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, imageURL, id)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *DamageModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM damages WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *DamageModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM damages WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *DamageModel) GetTotalDamagesByStore(ctx context.Context, storeID int64) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM damages WHERE store_id = ?`

	var total float64
	err := m.DB.QueryRowContext(ctx, query, storeID).Scan(&total)
	return total, err
}

func (m *DamageModel) GetTotalQuantityByStore(ctx context.Context, storeID int64) (float64, error) {
	query := `SELECT COALESCE(SUM(quantity), 0) FROM damages WHERE store_id = ?`

	var total float64
	err := m.DB.QueryRowContext(ctx, query, storeID).Scan(&total)
	return total, err
}

func (m *DamageModel) GetTotalWeightByStore(ctx context.Context, storeID int64) (float64, error) {
	query := `SELECT COALESCE(SUM(weight), 0) FROM damages WHERE store_id = ?`

	var total float64
	err := m.DB.QueryRowContext(ctx, query, storeID).Scan(&total)
	return total, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *DamageModel) scanDamage(row *sql.Row) (*Damage, error) {
	damage := &Damage{}
	err := row.Scan(
		&damage.ID,
		&damage.StoreID,
		&damage.UserID,
		&damage.Quantity,
		&damage.QuantityUnit,
		&damage.Weight,
		&damage.WeightUnit,
		&damage.DamageDate,
		&damage.Reason,
		&damage.Amount,
		&damage.Notes,
		&damage.ImageURL,
		&damage.CreatedAt,
		&damage.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return damage, nil
}

func (m *DamageModel) scanDamageRow(rows *sql.Rows) (*Damage, error) {
	damage := &Damage{}
	err := rows.Scan(
		&damage.ID,
		&damage.StoreID,
		&damage.UserID,
		&damage.Quantity,
		&damage.QuantityUnit,
		&damage.Weight,
		&damage.WeightUnit,
		&damage.DamageDate,
		&damage.Reason,
		&damage.Amount,
		&damage.Notes,
		&damage.ImageURL,
		&damage.CreatedAt,
		&damage.UpdatedAt,
	)
	return damage, err
}
