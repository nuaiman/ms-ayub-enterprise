package models

import (
	"context"
	"database/sql"
	"time"
)

type Transport struct {
	ID                     int64     `json:"id"`
	UserID                 int64     `json:"user_id"`
	CustomerID             *int64    `json:"customer_id,omitempty"`
	FromLocation           string    `json:"from_location"`
	ToLocation             *string   `json:"to_location,omitempty"`
	VehicleQuantity        float64   `json:"vehicle_quantity"`
	DeliveryType           *string   `json:"delivery_type,omitempty"` // local, district
	Notes                  *string   `json:"notes,omitempty"`
	TransportDate          time.Time `json:"transport_date"`
	OfficeCommissionAmount float64   `json:"office_commission_amount"`
	ImageURL               *string   `json:"image_url,omitempty"`
	CustomerTotalPaid      float64   `json:"customer_total_paid"` // NEW
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type TransportModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *TransportModel) Insert(ctx context.Context, transport *Transport) (int64, error) {
	query := `
		INSERT INTO transports (
			user_id, customer_id, from_location, to_location, vehicle_quantity,
			delivery_type, notes, transport_date, office_commission_amount, image_url,
			customer_total_paid
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		transport.UserID,
		transport.CustomerID,
		transport.FromLocation,
		transport.ToLocation,
		transport.VehicleQuantity,
		transport.DeliveryType,
		transport.Notes,
		transport.TransportDate,
		transport.OfficeCommissionAmount,
		transport.ImageURL,
		transport.CustomerTotalPaid, // NEW - defaults to 0
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *TransportModel) GetByID(ctx context.Context, id int64) (*Transport, error) {
	query := `
		SELECT id, user_id, customer_id, from_location, to_location, vehicle_quantity,
		       delivery_type, notes, transport_date, office_commission_amount, image_url,
		       customer_total_paid, created_at, updated_at
		FROM transports
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanTransport(row)
}

func (m *TransportModel) GetByCustomerID(ctx context.Context, customerID int64) ([]Transport, error) {
	query := `
		SELECT id, user_id, customer_id, from_location, to_location, vehicle_quantity,
		       delivery_type, notes, transport_date, office_commission_amount, image_url,
		       customer_total_paid, created_at, updated_at
		FROM transports
		WHERE customer_id = ?
		ORDER BY transport_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transports := []Transport{}
	for rows.Next() {
		transport, err := m.scanTransportRow(rows)
		if err != nil {
			return nil, err
		}
		transports = append(transports, *transport)
	}

	return transports, rows.Err()
}

func (m *TransportModel) GetByUserID(ctx context.Context, userID int64) ([]Transport, error) {
	query := `
		SELECT id, user_id, customer_id, from_location, to_location, vehicle_quantity,
		       delivery_type, notes, transport_date, office_commission_amount, image_url,
		       customer_total_paid, created_at, updated_at
		FROM transports
		WHERE user_id = ?
		ORDER BY transport_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transports := []Transport{}
	for rows.Next() {
		transport, err := m.scanTransportRow(rows)
		if err != nil {
			return nil, err
		}
		transports = append(transports, *transport)
	}

	return transports, rows.Err()
}

func (m *TransportModel) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]Transport, error) {
	query := `
		SELECT id, user_id, customer_id, from_location, to_location, vehicle_quantity,
		       delivery_type, notes, transport_date, office_commission_amount, image_url,
		       customer_total_paid, created_at, updated_at
		FROM transports
		WHERE transport_date >= ? AND transport_date <= ?
		ORDER BY transport_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transports := []Transport{}
	for rows.Next() {
		transport, err := m.scanTransportRow(rows)
		if err != nil {
			return nil, err
		}
		transports = append(transports, *transport)
	}

	return transports, rows.Err()
}

func (m *TransportModel) GetAll(ctx context.Context) ([]Transport, error) {
	query := `
		SELECT id, user_id, customer_id, from_location, to_location, vehicle_quantity,
		       delivery_type, notes, transport_date, office_commission_amount, image_url,
		       customer_total_paid, created_at, updated_at
		FROM transports
		ORDER BY transport_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transports := []Transport{}
	for rows.Next() {
		transport, err := m.scanTransportRow(rows)
		if err != nil {
			return nil, err
		}
		transports = append(transports, *transport)
	}

	return transports, rows.Err()
}

func (m *TransportModel) Search(ctx context.Context, query string) ([]Transport, error) {
	searchQuery := `
		SELECT id, user_id, customer_id, from_location, to_location, vehicle_quantity,
		       delivery_type, notes, transport_date, office_commission_amount, image_url,
		       customer_total_paid, created_at, updated_at
		FROM transports
		WHERE from_location LIKE ? OR to_location LIKE ?
		ORDER BY transport_date DESC
	`

	searchTerm := "%" + query + "%"

	rows, err := m.DB.QueryContext(ctx, searchQuery, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transports := []Transport{}
	for rows.Next() {
		transport, err := m.scanTransportRow(rows)
		if err != nil {
			return nil, err
		}
		transports = append(transports, *transport)
	}

	return transports, rows.Err()
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *TransportModel) Update(ctx context.Context, transport *Transport) error {
	query := `
		UPDATE transports
		SET 
			customer_id = ?,
			from_location = ?,
			to_location = ?,
			vehicle_quantity = ?,
			delivery_type = ?,
			notes = ?,
			transport_date = ?,
			office_commission_amount = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		transport.CustomerID,
		transport.FromLocation,
		transport.ToLocation,
		transport.VehicleQuantity,
		transport.DeliveryType,
		transport.Notes,
		transport.TransportDate,
		transport.OfficeCommissionAmount,
		transport.ID,
	)
	return err
}

func (m *TransportModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	query := `
		UPDATE transports
		SET image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, imageURL, id)
	return err
}

// NEW: Update customer total paid
func (m *TransportModel) UpdateCustomerPayment(ctx context.Context, id int64, totalPaid float64) error {
	query := `
		UPDATE transports
		SET 
			customer_total_paid = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, totalPaid, id)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *TransportModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM transports WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *TransportModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM transports WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *TransportModel) CountByCustomer(ctx context.Context, customerID int64) (int, error) {
	query := `SELECT COUNT(*) FROM transports WHERE customer_id = ?`

	var count int
	err := m.DB.QueryRowContext(ctx, query, customerID).Scan(&count)
	return count, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *TransportModel) scanTransport(row *sql.Row) (*Transport, error) {
	transport := &Transport{}
	err := row.Scan(
		&transport.ID,
		&transport.UserID,
		&transport.CustomerID,
		&transport.FromLocation,
		&transport.ToLocation,
		&transport.VehicleQuantity,
		&transport.DeliveryType,
		&transport.Notes,
		&transport.TransportDate,
		&transport.OfficeCommissionAmount,
		&transport.ImageURL,
		&transport.CustomerTotalPaid, // NEW
		&transport.CreatedAt,
		&transport.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return transport, nil
}

func (m *TransportModel) scanTransportRow(rows *sql.Rows) (*Transport, error) {
	transport := &Transport{}
	err := rows.Scan(
		&transport.ID,
		&transport.UserID,
		&transport.CustomerID,
		&transport.FromLocation,
		&transport.ToLocation,
		&transport.VehicleQuantity,
		&transport.DeliveryType,
		&transport.Notes,
		&transport.TransportDate,
		&transport.OfficeCommissionAmount,
		&transport.ImageURL,
		&transport.CustomerTotalPaid, // NEW
		&transport.CreatedAt,
		&transport.UpdatedAt,
	)
	return transport, err
}
