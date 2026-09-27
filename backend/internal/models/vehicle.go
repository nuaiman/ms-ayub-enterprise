package models

import (
	"context"
	"database/sql"
	"time"
)

type Vehicle struct {
	ID                int64     `json:"id"`
	UserID            int64     `json:"user_id"`
	TransportID       int64     `json:"transport_id"`
	VehicleNumber     string    `json:"vehicle_number"`
	BrokerID          int64     `json:"broker_id"`
	JomaCost          float64   `json:"joma_cost"`
	VehicleCost       float64   `json:"vehicle_cost"`
	TotalPaidToBroker float64   `json:"total_paid_to_broker"`
	OtherCost         float64   `json:"other_cost"`
	LabourCost        float64   `json:"labour_cost"`
	DemarageCost      float64   `json:"demarage_cost"`
	Notes             *string   `json:"notes,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
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
			user_id, transport_id, vehicle_number, broker_id,
			joma_cost, vehicle_cost, total_paid_to_broker,
			other_cost, labour_cost, demarage_cost, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		vehicle.UserID,
		vehicle.TransportID,
		vehicle.VehicleNumber,
		vehicle.BrokerID,
		vehicle.JomaCost,
		vehicle.VehicleCost,
		vehicle.TotalPaidToBroker,
		vehicle.OtherCost,
		vehicle.LabourCost,
		vehicle.DemarageCost,
		vehicle.Notes,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

const vehicleSelectCols = `
	id, user_id, transport_id, vehicle_number, broker_id,
	joma_cost, vehicle_cost, total_paid_to_broker,
	other_cost, labour_cost, demarage_cost, notes,
	created_at, updated_at
`

func (m *VehicleModel) GetByID(ctx context.Context, id int64) (*Vehicle, error) {
	query := `SELECT ` + vehicleSelectCols + ` FROM vehicles WHERE id = ?`
	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanVehicle(row)
}

func (m *VehicleModel) GetByTransportID(ctx context.Context, transportID int64) ([]Vehicle, error) {
	query := `SELECT ` + vehicleSelectCols + ` FROM vehicles WHERE transport_id = ? ORDER BY id ASC`
	rows, err := m.DB.QueryContext(ctx, query, transportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vehicles := []Vehicle{}
	for rows.Next() {
		v, err := m.scanVehicleRow(rows)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, *v)
	}
	return vehicles, rows.Err()
}

func (m *VehicleModel) GetByBrokerID(ctx context.Context, brokerID int64) ([]Vehicle, error) {
	query := `SELECT ` + vehicleSelectCols + ` FROM vehicles WHERE broker_id = ? ORDER BY created_at DESC`
	rows, err := m.DB.QueryContext(ctx, query, brokerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vehicles := []Vehicle{}
	for rows.Next() {
		v, err := m.scanVehicleRow(rows)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, *v)
	}
	return vehicles, rows.Err()
}

func (m *VehicleModel) GetByVehicleNumber(ctx context.Context, vehicleNumber string) ([]Vehicle, error) {
	query := `SELECT ` + vehicleSelectCols + ` FROM vehicles WHERE vehicle_number LIKE ? ORDER BY created_at DESC`
	searchTerm := "%" + vehicleNumber + "%"

	rows, err := m.DB.QueryContext(ctx, query, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vehicles := []Vehicle{}
	for rows.Next() {
		v, err := m.scanVehicleRow(rows)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, *v)
	}
	return vehicles, rows.Err()
}

func (m *VehicleModel) GetAll(ctx context.Context) ([]Vehicle, error) {
	query := `SELECT ` + vehicleSelectCols + ` FROM vehicles ORDER BY created_at DESC`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vehicles := []Vehicle{}
	for rows.Next() {
		v, err := m.scanVehicleRow(rows)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, *v)
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
			joma_cost = ?,
			vehicle_cost = ?,
			total_paid_to_broker = ?,
			other_cost = ?,
			labour_cost = ?,
			demarage_cost = ?,
			notes = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		vehicle.VehicleNumber,
		vehicle.BrokerID,
		vehicle.JomaCost,
		vehicle.VehicleCost,
		vehicle.TotalPaidToBroker,
		vehicle.OtherCost,
		vehicle.LabourCost,
		vehicle.DemarageCost,
		vehicle.Notes,
		vehicle.ID,
	)
	return err
}

func (m *VehicleModel) UpdateBrokerPayment(ctx context.Context, id int64, totalPaid float64) error {
	query := `
		UPDATE vehicles
		SET 
			total_paid_to_broker = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := m.DB.ExecContext(ctx, query, totalPaid, id)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *VehicleModel) Delete(ctx context.Context, id int64) error {
	_, err := m.DB.ExecContext(ctx, `DELETE FROM vehicles WHERE id = ?`, id)
	return err
}

func (m *VehicleModel) DeleteByTransportID(ctx context.Context, transportID int64) error {
	_, err := m.DB.ExecContext(ctx, `DELETE FROM vehicles WHERE transport_id = ?`, transportID)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *VehicleModel) Exists(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := m.DB.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM vehicles WHERE id = ?)`, id).Scan(&exists)
	return exists, err
}

func (m *VehicleModel) ExistsByNumber(ctx context.Context, vehicleNumber string) (bool, error) {
	var exists bool
	err := m.DB.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM vehicles WHERE vehicle_number = ?)`, vehicleNumber).Scan(&exists)
	return exists, err
}

func (m *VehicleModel) CountByTransport(ctx context.Context, transportID int64) (int, error) {
	var count int
	err := m.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM vehicles WHERE transport_id = ?`, transportID).Scan(&count)
	return count, err
}

func (m *VehicleModel) GetTotalCostByTransport(ctx context.Context, transportID int64) (float64, error) {
	query := `
		SELECT COALESCE(SUM(joma_cost + vehicle_cost + other_cost + labour_cost + demarage_cost), 0)
		FROM vehicles
		WHERE transport_id = ?
	`
	var total float64
	err := m.DB.QueryRowContext(ctx, query, transportID).Scan(&total)
	return total, err
}

func (m *VehicleModel) GetTotalJomaPlusVehicleByTransport(ctx context.Context, transportID int64) (float64, error) {
	query := `
		SELECT COALESCE(SUM(joma_cost + vehicle_cost), 0)
		FROM vehicles
		WHERE transport_id = ?
	`
	var total float64
	err := m.DB.QueryRowContext(ctx, query, transportID).Scan(&total)
	return total, err
}

func (m *VehicleModel) GetTotalPaidToBrokerByTransport(ctx context.Context, transportID int64) (float64, error) {
	query := `
		SELECT COALESCE(SUM(total_paid_to_broker), 0)
		FROM vehicles
		WHERE transport_id = ?
	`
	var total float64
	err := m.DB.QueryRowContext(ctx, query, transportID).Scan(&total)
	return total, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *VehicleModel) scanVehicle(row *sql.Row) (*Vehicle, error) {
	v := &Vehicle{}
	err := row.Scan(
		&v.ID,
		&v.UserID,
		&v.TransportID,
		&v.VehicleNumber,
		&v.BrokerID,
		&v.JomaCost,
		&v.VehicleCost,
		&v.TotalPaidToBroker,
		&v.OtherCost,
		&v.LabourCost,
		&v.DemarageCost,
		&v.Notes,
		&v.CreatedAt,
		&v.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (m *VehicleModel) scanVehicleRow(rows *sql.Rows) (*Vehicle, error) {
	v := &Vehicle{}
	err := rows.Scan(
		&v.ID,
		&v.UserID,
		&v.TransportID,
		&v.VehicleNumber,
		&v.BrokerID,
		&v.JomaCost,
		&v.VehicleCost,
		&v.TotalPaidToBroker,
		&v.OtherCost,
		&v.LabourCost,
		&v.DemarageCost,
		&v.Notes,
		&v.CreatedAt,
		&v.UpdatedAt,
	)
	return v, err
}
