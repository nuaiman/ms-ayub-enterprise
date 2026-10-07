package models

import (
	"context"
	"database/sql"
	"time"
)

// =============================================================================
// DELIVERY (Master)
// =============================================================================

type Delivery struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	CustomerID    *int64    `json:"customer_id,omitempty"`
	DeliveryDate  time.Time `json:"delivery_date"`
	ReceiverName  *string   `json:"receiver_name,omitempty"`
	ReceiverPhone *string   `json:"receiver_phone,omitempty"`
	FromLocation  *string   `json:"from_location,omitempty"`
	ToLocation    *string   `json:"to_location,omitempty"`
	Notes         *string   `json:"notes,omitempty"`
	ImageURL      *string   `json:"image_url,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

const deliverySelectCols = `
	id, user_id, customer_id, delivery_date,
	receiver_name, receiver_phone,
	from_location, to_location,
	notes, image_url,
	created_at, updated_at
`

type DeliveryModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *DeliveryModel) Insert(ctx context.Context, d *Delivery) (int64, error) {
	query := `
		INSERT INTO deliveries (
			user_id, customer_id, delivery_date,
			receiver_name, receiver_phone,
			from_location, to_location,
			notes, image_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		d.UserID,
		d.CustomerID,
		d.DeliveryDate,
		d.ReceiverName,
		d.ReceiverPhone,
		d.FromLocation,
		d.ToLocation,
		d.Notes,
		d.ImageURL,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *DeliveryModel) GetByID(ctx context.Context, id int64) (*Delivery, error) {
	query := `SELECT ` + deliverySelectCols + ` FROM deliveries WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scan(row)
}

func (m *DeliveryModel) GetAll(ctx context.Context) ([]Delivery, error) {
	query := `SELECT ` + deliverySelectCols + ` FROM deliveries ORDER BY delivery_date DESC, id DESC`
	return m.query(ctx, query)
}

func (m *DeliveryModel) GetByCustomerID(ctx context.Context, customerID int64) ([]Delivery, error) {
	query := `SELECT ` + deliverySelectCols + ` FROM deliveries WHERE customer_id = ? ORDER BY delivery_date DESC, id DESC`
	return m.query(ctx, query, customerID)
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *DeliveryModel) Update(ctx context.Context, d *Delivery) error {
	query := `
		UPDATE deliveries
		SET
			customer_id = ?,
			delivery_date = ?,
			receiver_name = ?,
			receiver_phone = ?,
			from_location = ?,
			to_location = ?,
			notes = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := m.DB.ExecContext(ctx, query,
		d.CustomerID,
		d.DeliveryDate,
		d.ReceiverName,
		d.ReceiverPhone,
		d.FromLocation,
		d.ToLocation,
		d.Notes,
		d.ID,
	)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *DeliveryModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM deliveries WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *DeliveryModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM deliveries WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS / HELPERS
// =============================================================================

func (m *DeliveryModel) query(ctx context.Context, query string, args ...any) ([]Delivery, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deliveries := []Delivery{}
	for rows.Next() {
		d, err := m.scanRow(rows)
		if err != nil {
			return nil, err
		}
		deliveries = append(deliveries, *d)
	}
	return deliveries, rows.Err()
}

func (m *DeliveryModel) scan(row *sql.Row) (*Delivery, error) {
	d := &Delivery{}
	err := row.Scan(
		&d.ID,
		&d.UserID,
		&d.CustomerID,
		&d.DeliveryDate,
		&d.ReceiverName,
		&d.ReceiverPhone,
		&d.FromLocation,
		&d.ToLocation,
		&d.Notes,
		&d.ImageURL,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return d, nil
}

func (m *DeliveryModel) scanRow(rows *sql.Rows) (*Delivery, error) {
	d := &Delivery{}
	err := rows.Scan(
		&d.ID,
		&d.UserID,
		&d.CustomerID,
		&d.DeliveryDate,
		&d.ReceiverName,
		&d.ReceiverPhone,
		&d.FromLocation,
		&d.ToLocation,
		&d.Notes,
		&d.ImageURL,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
	return d, err
}

// =============================================================================
// DELIVERY ITEM (Detail)
// =============================================================================

type DeliveryItem struct {
	ID            int64     `json:"id"`
	DeliveryID    int64     `json:"delivery_id"`
	StoreID       int64     `json:"store_id"`
	MajhiID       int64     `json:"majhi_id"`
	VehicleNumber *string   `json:"vehicle_number,omitempty"`
	DriverNumber  *string   `json:"driver_number,omitempty"`
	Quantity      float64   `json:"quantity"`
	Weight        float64   `json:"weight"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

const deliveryItemSelectCols = `
	id, delivery_id, store_id, majhi_id,
	vehicle_number, driver_number,
	quantity, weight,
	created_at, updated_at
`

// =============================================================================
// CREATE
// =============================================================================

func (m *DeliveryModel) InsertItem(ctx context.Context, it *DeliveryItem) (int64, error) {
	query := `
		INSERT INTO delivery_items (
			delivery_id, store_id, majhi_id,
			vehicle_number, driver_number,
			quantity, weight
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		it.DeliveryID,
		it.StoreID,
		it.MajhiID,
		it.VehicleNumber,
		it.DriverNumber,
		it.Quantity,
		it.Weight,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *DeliveryModel) GetItemByID(ctx context.Context, id int64) (*DeliveryItem, error) {
	query := `SELECT ` + deliveryItemSelectCols + ` FROM delivery_items WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanItem(row)
}

func (m *DeliveryModel) GetItemsByDeliveryID(ctx context.Context, deliveryID int64) ([]DeliveryItem, error) {
	query := `SELECT ` + deliveryItemSelectCols + ` FROM delivery_items WHERE delivery_id = ? ORDER BY id ASC`
	return m.queryItems(ctx, query, deliveryID)
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *DeliveryModel) UpdateItem(ctx context.Context, it *DeliveryItem) error {
	query := `
		UPDATE delivery_items
		SET
			store_id = ?,
			majhi_id = ?,
			vehicle_number = ?,
			driver_number = ?,
			quantity = ?,
			weight = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := m.DB.ExecContext(ctx, query,
		it.StoreID,
		it.MajhiID,
		it.VehicleNumber,
		it.DriverNumber,
		it.Quantity,
		it.Weight,
		it.ID,
	)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *DeliveryModel) DeleteItem(ctx context.Context, id int64) error {
	query := `DELETE FROM delivery_items WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *DeliveryModel) ItemExists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM delivery_items WHERE id = ?)`
	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// =============================================================================
// IMAGE
// =============================================================================

func (m *DeliveryModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	query := `
		UPDATE deliveries
		SET image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := m.DB.ExecContext(ctx, query, imageURL, id)
	return err
}

// =============================================================================
// SCANNERS / HELPERS
// =============================================================================

func (m *DeliveryModel) queryItems(ctx context.Context, query string, args ...any) ([]DeliveryItem, error) {
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []DeliveryItem{}
	for rows.Next() {
		it, err := m.scanItemRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *it)
	}
	return items, rows.Err()
}

func (m *DeliveryModel) scanItem(row *sql.Row) (*DeliveryItem, error) {
	it := &DeliveryItem{}
	err := row.Scan(
		&it.ID,
		&it.DeliveryID,
		&it.StoreID,
		&it.MajhiID,
		&it.VehicleNumber,
		&it.DriverNumber,
		&it.Quantity,
		&it.Weight,
		&it.CreatedAt,
		&it.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return it, nil
}

func (m *DeliveryModel) scanItemRow(rows *sql.Rows) (*DeliveryItem, error) {
	it := &DeliveryItem{}
	err := rows.Scan(
		&it.ID,
		&it.DeliveryID,
		&it.StoreID,
		&it.MajhiID,
		&it.VehicleNumber,
		&it.DriverNumber,
		&it.Quantity,
		&it.Weight,
		&it.CreatedAt,
		&it.UpdatedAt,
	)
	return it, err
}
