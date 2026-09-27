package models

import (
	"context"
	"database/sql"
	"time"
)

type Transport struct {
	ID                       int64      `json:"id"`
	UserID                   int64      `json:"user_id"`
	CustomerID               int64      `json:"customer_id"`
	FromLocation             string     `json:"from_location"`
	ToLocation               *string    `json:"to_location,omitempty"`
	VehicleQuantity          float64    `json:"vehicle_quantity"`
	TransportDate            time.Time  `json:"transport_date"`
	TransportType            *string    `json:"transport_type,omitempty"`
	ImageURL                 *string    `json:"image_url,omitempty"`
	Notes                    *string    `json:"notes,omitempty"`
	OfficeCommissionAmount   float64    `json:"office_commission_amount"`

	CustomerChargeUnit       string     `json:"customer_charge_unit"`
	CustomerTotalUnit        float64    `json:"customer_total_unit"`
	CustomerChargePerUnit    float64    `json:"customer_charge_per_unit"`
	CustomerTotalCharge      float64    `json:"customer_total_charge"`
	CustomerTotalPaid        float64    `json:"customer_total_paid"`
	CustomerTotalPaidThrough *time.Time `json:"customer_total_paid_through,omitempty"`

	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
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
			user_id, customer_id, from_location, to_location,
			vehicle_quantity, transport_date, transport_type,
			image_url, notes, office_commission_amount,
			customer_charge_unit, customer_total_unit, customer_charge_per_unit,
			customer_total_charge, customer_total_paid, customer_total_paid_through
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		transport.UserID,
		transport.CustomerID,
		transport.FromLocation,
		transport.ToLocation,
		transport.VehicleQuantity,
		transport.TransportDate,
		transport.TransportType,
		transport.ImageURL,
		transport.Notes,
		transport.OfficeCommissionAmount,
		transport.CustomerChargeUnit,
		transport.CustomerTotalUnit,
		transport.CustomerChargePerUnit,
		transport.CustomerTotalCharge,
		transport.CustomerTotalPaid,
		transport.CustomerTotalPaidThrough,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

const transportSelectCols = `
	id, user_id, customer_id, from_location, to_location,
	vehicle_quantity, transport_date, transport_type,
	image_url, notes, office_commission_amount,
	customer_charge_unit, customer_total_unit, customer_charge_per_unit,
	customer_total_charge, customer_total_paid, customer_total_paid_through,
	created_at, updated_at
`

func (m *TransportModel) GetByID(ctx context.Context, id int64) (*Transport, error) {
	query := `SELECT ` + transportSelectCols + ` FROM transports WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanTransport(row)
}

func (m *TransportModel) GetByCustomerID(ctx context.Context, customerID int64) ([]Transport, error) {
	query := `SELECT ` + transportSelectCols + ` FROM transports WHERE customer_id = ? ORDER BY transport_date DESC`
	rows, err := m.DB.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transports := []Transport{}
	for rows.Next() {
		t, err := m.scanTransportRow(rows)
		if err != nil {
			return nil, err
		}
		transports = append(transports, *t)
	}
	return transports, rows.Err()
}

func (m *TransportModel) GetByUserID(ctx context.Context, userID int64) ([]Transport, error) {
	query := `SELECT ` + transportSelectCols + ` FROM transports WHERE user_id = ? ORDER BY transport_date DESC`
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transports := []Transport{}
	for rows.Next() {
		t, err := m.scanTransportRow(rows)
		if err != nil {
			return nil, err
		}
		transports = append(transports, *t)
	}
	return transports, rows.Err()
}

func (m *TransportModel) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]Transport, error) {
	query := `SELECT ` + transportSelectCols + ` FROM transports WHERE transport_date >= ? AND transport_date <= ? ORDER BY transport_date DESC`
	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transports := []Transport{}
	for rows.Next() {
		t, err := m.scanTransportRow(rows)
		if err != nil {
			return nil, err
		}
		transports = append(transports, *t)
	}
	return transports, rows.Err()
}

func (m *TransportModel) GetAll(ctx context.Context) ([]Transport, error) {
	query := `SELECT ` + transportSelectCols + ` FROM transports ORDER BY transport_date DESC`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transports := []Transport{}
	for rows.Next() {
		t, err := m.scanTransportRow(rows)
		if err != nil {
			return nil, err
		}
		transports = append(transports, *t)
	}
	return transports, rows.Err()
}

func (m *TransportModel) Search(ctx context.Context, query string) ([]Transport, error) {
	searchQuery := `SELECT ` + transportSelectCols + ` FROM transports WHERE from_location LIKE ? OR to_location LIKE ? ORDER BY transport_date DESC`
	searchTerm := "%" + query + "%"

	rows, err := m.DB.QueryContext(ctx, searchQuery, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transports := []Transport{}
	for rows.Next() {
		t, err := m.scanTransportRow(rows)
		if err != nil {
			return nil, err
		}
		transports = append(transports, *t)
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
			transport_date = ?,
			transport_type = ?,
			notes = ?,
			office_commission_amount = ?,
			customer_charge_unit = ?,
			customer_total_unit = ?,
			customer_charge_per_unit = ?,
			customer_total_charge = ?,
			customer_total_paid = ?,
			customer_total_paid_through = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		transport.CustomerID,
		transport.FromLocation,
		transport.ToLocation,
		transport.VehicleQuantity,
		transport.TransportDate,
		transport.TransportType,
		transport.Notes,
		transport.OfficeCommissionAmount,
		transport.CustomerChargeUnit,
		transport.CustomerTotalUnit,
		transport.CustomerChargePerUnit,
		transport.CustomerTotalCharge,
		transport.CustomerTotalPaid,
		transport.CustomerTotalPaidThrough,
		transport.ID,
	)
	return err
}

func (m *TransportModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	query := `UPDATE transports SET image_url = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, imageURL, id)
	return err
}

func (m *TransportModel) UpdateCustomerPayment(ctx context.Context, id int64, totalPaid float64, paidThrough *time.Time) error {
	query := `
		UPDATE transports
		SET 
			customer_total_paid = ?,
			customer_total_paid_through = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := m.DB.ExecContext(ctx, query, totalPaid, paidThrough, id)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *TransportModel) Delete(ctx context.Context, id int64) error {
	_, err := m.DB.ExecContext(ctx, `DELETE FROM transports WHERE id = ?`, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *TransportModel) Exists(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := m.DB.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM transports WHERE id = ?)`, id).Scan(&exists)
	return exists, err
}

func (m *TransportModel) CountByCustomer(ctx context.Context, customerID int64) (int, error) {
	var count int
	err := m.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM transports WHERE customer_id = ?`, customerID).Scan(&count)
	return count, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *TransportModel) scanTransport(row *sql.Row) (*Transport, error) {
	t := &Transport{}
	err := row.Scan(
		&t.ID,
		&t.UserID,
		&t.CustomerID,
		&t.FromLocation,
		&t.ToLocation,
		&t.VehicleQuantity,
		&t.TransportDate,
		&t.TransportType,
		&t.ImageURL,
		&t.Notes,
		&t.OfficeCommissionAmount,
		&t.CustomerChargeUnit,
		&t.CustomerTotalUnit,
		&t.CustomerChargePerUnit,
		&t.CustomerTotalCharge,
		&t.CustomerTotalPaid,
		&t.CustomerTotalPaidThrough,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}

func (m *TransportModel) scanTransportRow(rows *sql.Rows) (*Transport, error) {
	t := &Transport{}
	err := rows.Scan(
		&t.ID,
		&t.UserID,
		&t.CustomerID,
		&t.FromLocation,
		&t.ToLocation,
		&t.VehicleQuantity,
		&t.TransportDate,
		&t.TransportType,
		&t.ImageURL,
		&t.Notes,
		&t.OfficeCommissionAmount,
		&t.CustomerChargeUnit,
		&t.CustomerTotalUnit,
		&t.CustomerChargePerUnit,
		&t.CustomerTotalCharge,
		&t.CustomerTotalPaid,
		&t.CustomerTotalPaidThrough,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	return t, err
}