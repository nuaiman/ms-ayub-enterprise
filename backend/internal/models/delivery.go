package models

import (
	"context"
	"database/sql"
	"time"
)

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

type DeliveryModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *DeliveryModel) Insert(ctx context.Context, delivery *Delivery) (int64, error) {
	query := `
		INSERT INTO deliveries (
			user_id, customer_id, delivery_date, receiver_name, receiver_phone,
			from_location, to_location, notes, image_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		delivery.UserID,
		delivery.CustomerID,
		delivery.DeliveryDate,
		delivery.ReceiverName,
		delivery.ReceiverPhone,
		delivery.FromLocation,
		delivery.ToLocation,
		delivery.Notes,
		delivery.ImageURL,
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
	query := `
		SELECT id, user_id, customer_id, delivery_date, receiver_name, receiver_phone,
		       from_location, to_location, notes, image_url, created_at, updated_at
		FROM deliveries
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanDelivery(row)
}

func (m *DeliveryModel) GetByCustomerID(ctx context.Context, customerID int64) ([]Delivery, error) {
	query := `
		SELECT id, user_id, customer_id, delivery_date, receiver_name, receiver_phone,
		       from_location, to_location, notes, image_url, created_at, updated_at
		FROM deliveries
		WHERE customer_id = ?
		ORDER BY delivery_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deliveries := []Delivery{}
	for rows.Next() {
		delivery, err := m.scanDeliveryRow(rows)
		if err != nil {
			return nil, err
		}
		deliveries = append(deliveries, *delivery)
	}

	return deliveries, rows.Err()
}

func (m *DeliveryModel) GetByUserID(ctx context.Context, userID int64) ([]Delivery, error) {
	query := `
		SELECT id, user_id, customer_id, delivery_date, receiver_name, receiver_phone,
		       from_location, to_location, notes, image_url, created_at, updated_at
		FROM deliveries
		WHERE user_id = ?
		ORDER BY delivery_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deliveries := []Delivery{}
	for rows.Next() {
		delivery, err := m.scanDeliveryRow(rows)
		if err != nil {
			return nil, err
		}
		deliveries = append(deliveries, *delivery)
	}

	return deliveries, rows.Err()
}

func (m *DeliveryModel) GetAll(ctx context.Context) ([]Delivery, error) {
	query := `
		SELECT id, user_id, customer_id, delivery_date, receiver_name, receiver_phone,
		       from_location, to_location, notes, image_url, created_at, updated_at
		FROM deliveries
		ORDER BY delivery_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deliveries := []Delivery{}
	for rows.Next() {
		delivery, err := m.scanDeliveryRow(rows)
		if err != nil {
			return nil, err
		}
		deliveries = append(deliveries, *delivery)
	}

	return deliveries, rows.Err()
}

func (m *DeliveryModel) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]Delivery, error) {
	query := `
		SELECT id, user_id, customer_id, delivery_date, receiver_name, receiver_phone,
		       from_location, to_location, notes, image_url, created_at, updated_at
		FROM deliveries
		WHERE delivery_date >= ? AND delivery_date <= ?
		ORDER BY delivery_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deliveries := []Delivery{}
	for rows.Next() {
		delivery, err := m.scanDeliveryRow(rows)
		if err != nil {
			return nil, err
		}
		deliveries = append(deliveries, *delivery)
	}

	return deliveries, rows.Err()
}

func (m *DeliveryModel) Search(ctx context.Context, query string) ([]Delivery, error) {
	searchQuery := `
		SELECT id, user_id, customer_id, delivery_date, receiver_name, receiver_phone,
		       from_location, to_location, notes, image_url, created_at, updated_at
		FROM deliveries
		WHERE receiver_name LIKE ? OR receiver_phone LIKE ? OR from_location LIKE ? OR to_location LIKE ?
		ORDER BY delivery_date DESC
	`

	searchTerm := "%" + query + "%"

	rows, err := m.DB.QueryContext(ctx, searchQuery,
		searchTerm, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deliveries := []Delivery{}
	for rows.Next() {
		delivery, err := m.scanDeliveryRow(rows)
		if err != nil {
			return nil, err
		}
		deliveries = append(deliveries, *delivery)
	}

	return deliveries, rows.Err()
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *DeliveryModel) Update(ctx context.Context, delivery *Delivery) error {
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
		delivery.CustomerID,
		delivery.DeliveryDate,
		delivery.ReceiverName,
		delivery.ReceiverPhone,
		delivery.FromLocation,
		delivery.ToLocation,
		delivery.Notes,
		delivery.ID,
	)
	return err
}

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

func (m *DeliveryModel) CountByCustomer(ctx context.Context, customerID int64) (int, error) {
	query := `SELECT COUNT(*) FROM deliveries WHERE customer_id = ?`

	var count int
	err := m.DB.QueryRowContext(ctx, query, customerID).Scan(&count)
	return count, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *DeliveryModel) scanDelivery(row *sql.Row) (*Delivery, error) {
	delivery := &Delivery{}
	err := row.Scan(
		&delivery.ID,
		&delivery.UserID,
		&delivery.CustomerID,
		&delivery.DeliveryDate,
		&delivery.ReceiverName,
		&delivery.ReceiverPhone,
		&delivery.FromLocation,
		&delivery.ToLocation,
		&delivery.Notes,
		&delivery.ImageURL,
		&delivery.CreatedAt,
		&delivery.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return delivery, nil
}

func (m *DeliveryModel) scanDeliveryRow(rows *sql.Rows) (*Delivery, error) {
	delivery := &Delivery{}
	err := rows.Scan(
		&delivery.ID,
		&delivery.UserID,
		&delivery.CustomerID,
		&delivery.DeliveryDate,
		&delivery.ReceiverName,
		&delivery.ReceiverPhone,
		&delivery.FromLocation,
		&delivery.ToLocation,
		&delivery.Notes,
		&delivery.ImageURL,
		&delivery.CreatedAt,
		&delivery.UpdatedAt,
	)
	return delivery, err
}
