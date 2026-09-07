package models

import (
	"context"
	"database/sql"
	"time"
)

type Broker struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone,omitempty"`
	Notes     *string   `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BrokerModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *BrokerModel) Insert(ctx context.Context, broker *Broker) (int64, error) {
	query := `
		INSERT INTO brokers (name, phone, notes)
		VALUES (?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		broker.Name,
		broker.Phone,
		broker.Notes,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *BrokerModel) GetByID(ctx context.Context, id int64) (*Broker, error) {
	query := `
		SELECT id, name, phone, notes, created_at, updated_at
		FROM brokers
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanBroker(row)
}

func (m *BrokerModel) GetAll(ctx context.Context) ([]Broker, error) {
	query := `
		SELECT id, name, phone, notes, created_at, updated_at
		FROM brokers
		ORDER BY name ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	brokers := []Broker{}
	for rows.Next() {
		broker, err := m.scanBrokerRow(rows)
		if err != nil {
			return nil, err
		}
		brokers = append(brokers, *broker)
	}

	return brokers, rows.Err()
}

func (m *BrokerModel) Search(ctx context.Context, query string) ([]Broker, error) {
	searchQuery := `
		SELECT id, name, phone, notes, created_at, updated_at
		FROM brokers
		WHERE name LIKE ? OR phone LIKE ?
		ORDER BY name ASC
	`

	searchTerm := "%" + query + "%"

	rows, err := m.DB.QueryContext(ctx, searchQuery, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	brokers := []Broker{}
	for rows.Next() {
		broker, err := m.scanBrokerRow(rows)
		if err != nil {
			return nil, err
		}
		brokers = append(brokers, *broker)
	}

	return brokers, rows.Err()
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *BrokerModel) Update(ctx context.Context, broker *Broker) error {
	query := `
		UPDATE brokers
		SET name = ?, phone = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		broker.Name,
		broker.Phone,
		broker.Notes,
		broker.ID,
	)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *BrokerModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM brokers WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *BrokerModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM brokers WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *BrokerModel) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM brokers WHERE name = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, name).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *BrokerModel) scanBroker(row *sql.Row) (*Broker, error) {
	broker := &Broker{}
	err := row.Scan(
		&broker.ID,
		&broker.Name,
		&broker.Phone,
		&broker.Notes,
		&broker.CreatedAt,
		&broker.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return broker, nil
}

func (m *BrokerModel) scanBrokerRow(rows *sql.Rows) (*Broker, error) {
	broker := &Broker{}
	err := rows.Scan(
		&broker.ID,
		&broker.Name,
		&broker.Phone,
		&broker.Notes,
		&broker.CreatedAt,
		&broker.UpdatedAt,
	)
	return broker, err
}
