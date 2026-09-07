package models

import (
	"context"
	"database/sql"
	"time"
)

type Customer struct {
	ID            int64     `json:"id"`
	CompanyName   *string   `json:"company_name,omitempty"`
	ContactPerson *string   `json:"contact_person,omitempty"`
	Phone         string    `json:"phone"`
	Email         *string   `json:"email,omitempty"`
	Address       *string   `json:"address,omitempty"`
	Notes         *string   `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CustomerModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *CustomerModel) Insert(ctx context.Context, customer *Customer) (int64, error) {
	query := `
		INSERT INTO customers (company_name, contact_person, phone, email, address, notes)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		customer.CompanyName,
		customer.ContactPerson,
		customer.Phone,
		customer.Email,
		customer.Address,
		customer.Notes,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *CustomerModel) GetByID(ctx context.Context, id int64) (*Customer, error) {
	query := `
		SELECT id, company_name, contact_person, phone, email, address, notes, created_at, updated_at
		FROM customers
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanCustomer(row)
}

func (m *CustomerModel) GetByPhone(ctx context.Context, phone string) (*Customer, error) {
	query := `
		SELECT id, company_name, contact_person, phone, email, address, notes, created_at, updated_at
		FROM customers
		WHERE phone = ?
	`

	row := m.DB.QueryRowContext(ctx, query, phone)
	return m.scanCustomer(row)
}

func (m *CustomerModel) GetAll(ctx context.Context) ([]Customer, error) {
	query := `
		SELECT id, company_name, contact_person, phone, email, address, notes, created_at, updated_at
		FROM customers
		ORDER BY company_name ASC, contact_person ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []Customer{}
	for rows.Next() {
		customer, err := m.scanCustomerRow(rows)
		if err != nil {
			return nil, err
		}
		customers = append(customers, *customer)
	}

	return customers, rows.Err()
}

func (m *CustomerModel) Search(ctx context.Context, query string) ([]Customer, error) {
	searchQuery := `
		SELECT id, company_name, contact_person, phone, email, address, notes, created_at, updated_at
		FROM customers
		WHERE company_name LIKE ? OR contact_person LIKE ? OR phone LIKE ? OR email LIKE ? OR address LIKE ?
		ORDER BY company_name ASC, contact_person ASC
	`

	searchTerm := "%" + query + "%"

	rows, err := m.DB.QueryContext(ctx, searchQuery,
		searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []Customer{}
	for rows.Next() {
		customer, err := m.scanCustomerRow(rows)
		if err != nil {
			return nil, err
		}
		customers = append(customers, *customer)
	}

	return customers, rows.Err()
}

func (m *CustomerModel) GetActiveCustomers(ctx context.Context) ([]Customer, error) {
	// Since customers table doesn't have is_active field, we just return all
	return m.GetAll(ctx)
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *CustomerModel) Update(ctx context.Context, customer *Customer) error {
	query := `
		UPDATE customers
		SET 
			company_name = ?,
			contact_person = ?,
			phone = ?,
			email = ?,
			address = ?,
			notes = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		customer.CompanyName,
		customer.ContactPerson,
		customer.Phone,
		customer.Email,
		customer.Address,
		customer.Notes,
		customer.ID,
	)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *CustomerModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM customers WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *CustomerModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customers WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *CustomerModel) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customers WHERE phone = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, phone).Scan(&exists)
	return exists, err
}

func (m *CustomerModel) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM customers WHERE email = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *CustomerModel) scanCustomer(row *sql.Row) (*Customer, error) {
	customer := &Customer{}
	err := row.Scan(
		&customer.ID,
		&customer.CompanyName,
		&customer.ContactPerson,
		&customer.Phone,
		&customer.Email,
		&customer.Address,
		&customer.Notes,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return customer, nil
}

func (m *CustomerModel) scanCustomerRow(rows *sql.Rows) (*Customer, error) {
	customer := &Customer{}
	err := rows.Scan(
		&customer.ID,
		&customer.CompanyName,
		&customer.ContactPerson,
		&customer.Phone,
		&customer.Email,
		&customer.Address,
		&customer.Notes,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	return customer, err
}
