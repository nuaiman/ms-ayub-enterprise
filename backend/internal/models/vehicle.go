package models

import (
	"context"
	"database/sql"
	"time"
)

type Vehicle struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	TransportID     int64     `json:"transport_id"`
	VehicleNumber   string    `json:"vehicle_number"`
	BrokerID        *int64    `json:"broker_id,omitempty"`
	DriverName      *string   `json:"driver_name,omitempty"`
	DriverPhone     *string   `json:"driver_phone,omitempty"`
	JomaCost        float64   `json:"joma_cost"`
	VehicleCost     float64   `json:"vehicle_cost"`
	CustomerCharge  float64   `json:"customer_charge"`
	OtherCost       float64   `json:"other_cost"`
	LabourCost      float64   `json:"labour_cost"`
	DemarageAmount  float64   `json:"demarage_amount"`
	DemarageReason  *string   `json:"demarage_reason,omitempty"`
	BrokerTotalPaid float64   `json:"broker_total_paid"` // NEW
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type VehicleModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *VehicleModel) Insert(ctx context.Context, vehicle *Vehicle) (int64, error) {
	query := `
		INSERT INTO vehicles (
			user_id, transport_id, vehicle_number, broker_id, driver_name, driver_phone,
			joma_cost, vehicle_cost, customer_charge, other_cost, labour_cost,
			demarage_amount, demarage_reason, broker_total_paid
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		vehicle.UserID,
		vehicle.TransportID,
		vehicle.VehicleNumber,
		vehicle.BrokerID,
		vehicle.DriverName,
		vehicle.DriverPhone,
		vehicle.JomaCost,
		vehicle.VehicleCost,
		vehicle.CustomerCharge,
		vehicle.OtherCost,
		vehicle.LabourCost,
		vehicle.DemarageAmount,
		vehicle.DemarageReason,
		vehicle.BrokerTotalPaid, // NEW
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *VehicleModel) GetByID(ctx context.Context, id int64) (*Vehicle, error) {
	query := `
		SELECT id, user_id, transport_id, vehicle_number, broker_id, driver_name, driver_phone,
		       joma_cost, vehicle_cost, customer_charge, other_cost, labour_cost,
		       demarage_amount, demarage_reason, broker_total_paid, created_at, updated_at
		FROM vehicles
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanVehicle(row)
}

func (m *VehicleModel) GetByTransportID(ctx context.Context, transportID int64) ([]Vehicle, error) {
	query := `
		SELECT id, user_id, transport_id, vehicle_number, broker_id, driver_name, driver_phone,
		       joma_cost, vehicle_cost, customer_charge, other_cost, labour_cost,
		       demarage_amount, demarage_reason, broker_total_paid, created_at, updated_at
		FROM vehicles
		WHERE transport_id = ?
		ORDER BY id ASC
	`

	rows, err := m.DB.QueryContext(ctx, query, transportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vehicles := []Vehicle{}
	for rows.Next() {
		vehicle, err := m.scanVehicleRow(rows)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, *vehicle)
	}

	return vehicles, rows.Err()
}

func (m *VehicleModel) GetByBrokerID(ctx context.Context, brokerID int64) ([]Vehicle, error) {
	query := `
		SELECT id, user_id, transport_id, vehicle_number, broker_id, driver_name, driver_phone,
		       joma_cost, vehicle_cost, customer_charge, other_cost, labour_cost,
		       demarage_amount, demarage_reason, broker_total_paid, created_at, updated_at
		FROM vehicles
		WHERE broker_id = ?
		ORDER BY created_at DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, brokerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vehicles := []Vehicle{}
	for rows.Next() {
		vehicle, err := m.scanVehicleRow(rows)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, *vehicle)
	}

	return vehicles, rows.Err()
}

func (m *VehicleModel) GetByVehicleNumber(ctx context.Context, vehicleNumber string) ([]Vehicle, error) {
	query := `
		SELECT id, user_id, transport_id, vehicle_number, broker_id, driver_name, driver_phone,
		       joma_cost, vehicle_cost, customer_charge, other_cost, labour_cost,
		       demarage_amount, demarage_reason, broker_total_paid, created_at, updated_at
		FROM vehicles
		WHERE vehicle_number LIKE ?
		ORDER BY created_at DESC
	`

	searchTerm := "%" + vehicleNumber + "%"

	rows, err := m.DB.QueryContext(ctx, query, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vehicles := []Vehicle{}
	for rows.Next() {
		vehicle, err := m.scanVehicleRow(rows)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, *vehicle)
	}

	return vehicles, rows.Err()
}

func (m *VehicleModel) GetAll(ctx context.Context) ([]Vehicle, error) {
	query := `
		SELECT id, user_id, transport_id, vehicle_number, broker_id, driver_name, driver_phone,
		       joma_cost, vehicle_cost, customer_charge, other_cost, labour_cost,
		       demarage_amount, demarage_reason, broker_total_paid, created_at, updated_at
		FROM vehicles
		ORDER BY created_at DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vehicles := []Vehicle{}
	for rows.Next() {
		vehicle, err := m.scanVehicleRow(rows)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, *vehicle)
	}

	return vehicles, rows.Err()
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *VehicleModel) Update(ctx context.Context, vehicle *Vehicle) error {
	query := `
		UPDATE vehicles
		SET 
			vehicle_number = ?,
			broker_id = ?,
			driver_name = ?,
			driver_phone = ?,
			joma_cost = ?,
			vehicle_cost = ?,
			customer_charge = ?,
			other_cost = ?,
			labour_cost = ?,
			demarage_amount = ?,
			demarage_reason = ?,
			broker_total_paid = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		vehicle.VehicleNumber,
		vehicle.BrokerID,
		vehicle.DriverName,
		vehicle.DriverPhone,
		vehicle.JomaCost,
		vehicle.VehicleCost,
		vehicle.CustomerCharge,
		vehicle.OtherCost,
		vehicle.LabourCost,
		vehicle.DemarageAmount,
		vehicle.DemarageReason,
		vehicle.BrokerTotalPaid, // NEW
		vehicle.ID,
	)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *VehicleModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM vehicles WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m *VehicleModel) DeleteByTransportID(ctx context.Context, transportID int64) error {
	query := `DELETE FROM vehicles WHERE transport_id = ?`
	_, err := m.DB.ExecContext(ctx, query, transportID)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *VehicleModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicles WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *VehicleModel) ExistsByNumber(ctx context.Context, vehicleNumber string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM vehicles WHERE vehicle_number = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, vehicleNumber).Scan(&exists)
	return exists, err
}

func (m *VehicleModel) CountByTransport(ctx context.Context, transportID int64) (int, error) {
	query := `SELECT COUNT(*) FROM vehicles WHERE transport_id = ?`

	var count int
	err := m.DB.QueryRowContext(ctx, query, transportID).Scan(&count)
	return count, err
}

func (m *VehicleModel) GetTotalCostByTransport(ctx context.Context, transportID int64) (float64, error) {
	query := `
		SELECT COALESCE(SUM(joma_cost + vehicle_cost + other_cost + labour_cost + demarage_amount), 0)
		FROM vehicles
		WHERE transport_id = ?
	`

	var total float64
	err := m.DB.QueryRowContext(ctx, query, transportID).Scan(&total)
	return total, err
}

// NEW: Update broker total paid
func (m *VehicleModel) UpdateBrokerPayment(ctx context.Context, id int64, totalPaid float64) error {
	query := `
		UPDATE vehicles
		SET 
			broker_total_paid = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, totalPaid, id)
	return err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *VehicleModel) scanVehicle(row *sql.Row) (*Vehicle, error) {
	vehicle := &Vehicle{}
	err := row.Scan(
		&vehicle.ID,
		&vehicle.UserID,
		&vehicle.TransportID,
		&vehicle.VehicleNumber,
		&vehicle.BrokerID,
		&vehicle.DriverName,
		&vehicle.DriverPhone,
		&vehicle.JomaCost,
		&vehicle.VehicleCost,
		&vehicle.CustomerCharge,
		&vehicle.OtherCost,
		&vehicle.LabourCost,
		&vehicle.DemarageAmount,
		&vehicle.DemarageReason,
		&vehicle.BrokerTotalPaid, // NEW
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return vehicle, nil
}

func (m *VehicleModel) scanVehicleRow(rows *sql.Rows) (*Vehicle, error) {
	vehicle := &Vehicle{}
	err := rows.Scan(
		&vehicle.ID,
		&vehicle.UserID,
		&vehicle.TransportID,
		&vehicle.VehicleNumber,
		&vehicle.BrokerID,
		&vehicle.DriverName,
		&vehicle.DriverPhone,
		&vehicle.JomaCost,
		&vehicle.VehicleCost,
		&vehicle.CustomerCharge,
		&vehicle.OtherCost,
		&vehicle.LabourCost,
		&vehicle.DemarageAmount,
		&vehicle.DemarageReason,
		&vehicle.BrokerTotalPaid, // NEW
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)
	return vehicle, err
}
