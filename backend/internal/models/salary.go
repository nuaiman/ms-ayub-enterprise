package models

import (
	"context"
	"database/sql"
	"time"
)

type Salary struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"user_id"`     // Created by
	EmployeeID      int64      `json:"employee_id"` // Employee being paid
	MonthYear       string     `json:"month_year"`  // Format: YYYY-MM
	Bonus           float64    `json:"bonus"`
	Deductions      float64    `json:"deductions"`
	Status          string     `json:"status"` // draft, paid, cancelled
	PaymentDate     *time.Time `json:"payment_date,omitempty"`
	PaymentMethod   *string    `json:"payment_method,omitempty"`
	ReferenceNumber *string    `json:"reference_number,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type SalaryModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *SalaryModel) Insert(ctx context.Context, salary *Salary) (int64, error) {
	query := `
		INSERT INTO salaries (
			user_id, employee_id, month_year, bonus, deductions, status, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		salary.UserID,
		salary.EmployeeID,
		salary.MonthYear,
		salary.Bonus,
		salary.Deductions,
		"draft",
		salary.Notes,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *SalaryModel) GetByID(ctx context.Context, id int64) (*Salary, error) {
	query := `
		SELECT 
			id, user_id, employee_id, month_year, 
			bonus, deductions,
			status, payment_date, payment_method, 
			reference_number, notes,
			created_at, updated_at
		FROM salaries
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanSalary(row)
}

func (m *SalaryModel) GetAll(ctx context.Context, filters map[string]string) ([]Salary, error) {
	query := `
		SELECT 
			id, user_id, employee_id, month_year, 
			bonus, deductions,
			status, payment_date, payment_method, 
			reference_number, notes,
			created_at, updated_at
		FROM salaries
		WHERE 1=1
	`

	args := []interface{}{}

	if employeeID, ok := filters["employee_id"]; ok && employeeID != "" {
		query += " AND employee_id = ?"
		args = append(args, employeeID)
	}

	if monthYear, ok := filters["month_year"]; ok && monthYear != "" {
		query += " AND month_year = ?"
		args = append(args, monthYear)
	}

	if status, ok := filters["status"]; ok && status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY month_year DESC, created_at DESC"

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	salaries := []Salary{}
	for rows.Next() {
		salary, err := m.scanSalaryRow(rows)
		if err != nil {
			return nil, err
		}
		salaries = append(salaries, *salary)
	}

	return salaries, rows.Err()
}

func (m *SalaryModel) GetByEmployeeID(ctx context.Context, employeeID int64) ([]Salary, error) {
	query := `
		SELECT 
			id, user_id, employee_id, month_year, 
			bonus, deductions,
			status, payment_date, payment_method, 
			reference_number, notes,
			created_at, updated_at
		FROM salaries
		WHERE employee_id = ?
		ORDER BY month_year DESC, created_at DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	salaries := []Salary{}
	for rows.Next() {
		salary, err := m.scanSalaryRow(rows)
		if err != nil {
			return nil, err
		}
		salaries = append(salaries, *salary)
	}

	return salaries, rows.Err()
}

func (m *SalaryModel) GetByMonth(ctx context.Context, monthYear string) ([]Salary, error) {
	filters := map[string]string{"month_year": monthYear}
	return m.GetAll(ctx, filters)
}

func (m *SalaryModel) GetByEmployeeAndMonth(ctx context.Context, employeeID int64, monthYear string) (*Salary, error) {
	query := `
		SELECT 
			id, user_id, employee_id, month_year, 
			bonus, deductions,
			status, payment_date, payment_method, 
			reference_number, notes,
			created_at, updated_at
		FROM salaries
		WHERE employee_id = ? AND month_year = ?
	`

	row := m.DB.QueryRowContext(ctx, query, employeeID, monthYear)
	return m.scanSalary(row)
}

func (m *SalaryModel) EnsureCurrentMonthSalaries(ctx context.Context, userID int64) error {
	currentMonth := time.Now().Format("2006-01")

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	userModel := UserModel{DB: m.DB}
	employees, err := userModel.GetActiveUsers(ctx)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT OR IGNORE INTO salaries (
			user_id, employee_id, month_year, bonus, deductions, status, notes
		) VALUES (?, ?, ?, 0, 0, 'draft', NULL)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, emp := range employees {
		if _, err := stmt.ExecContext(ctx, userID, emp.ID, currentMonth); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *SalaryModel) Update(ctx context.Context, salary *Salary) error {
	query := `
		UPDATE salaries
		SET 
			bonus = ?,
			deductions = ?,
			notes = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND status = 'draft'
	`

	result, err := m.DB.ExecContext(ctx, query,
		salary.Bonus,
		salary.Deductions,
		salary.Notes,
		salary.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (m *SalaryModel) MarkAsPaid(ctx context.Context, id int64, paymentDate *time.Time, paymentMethod, referenceNumber *string) error {
	query := `
		UPDATE salaries
		SET 
			payment_date = ?,
			payment_method = ?,
			reference_number = ?,
			status = 'paid',
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND status = 'draft'
	`

	result, err := m.DB.ExecContext(ctx, query,
		paymentDate,
		paymentMethod,
		referenceNumber,
		id,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (m *SalaryModel) Cancel(ctx context.Context, id int64) error {
	query := `
		UPDATE salaries
		SET 
			status = 'cancelled',
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND status IN ('draft', 'paid')
	`

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// =============================================================================
// DELETE
// =============================================================================

func (m *SalaryModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM salaries WHERE id = ? AND status = 'draft'`

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *SalaryModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM salaries WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *SalaryModel) ExistsForEmployeeMonth(ctx context.Context, employeeID int64, monthYear string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM salaries WHERE employee_id = ? AND month_year = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, employeeID, monthYear).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *SalaryModel) scanSalary(row *sql.Row) (*Salary, error) {
	salary := &Salary{}
	err := row.Scan(
		&salary.ID,
		&salary.UserID,
		&salary.EmployeeID,
		&salary.MonthYear,
		&salary.Bonus,
		&salary.Deductions,
		&salary.Status,
		&salary.PaymentDate,
		&salary.PaymentMethod,
		&salary.ReferenceNumber,
		&salary.Notes,
		&salary.CreatedAt,
		&salary.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return salary, nil
}

func (m *SalaryModel) scanSalaryRow(rows *sql.Rows) (*Salary, error) {
	salary := &Salary{}
	err := rows.Scan(
		&salary.ID,
		&salary.UserID,
		&salary.EmployeeID,
		&salary.MonthYear,
		&salary.Bonus,
		&salary.Deductions,
		&salary.Status,
		&salary.PaymentDate,
		&salary.PaymentMethod,
		&salary.ReferenceNumber,
		&salary.Notes,
		&salary.CreatedAt,
		&salary.UpdatedAt,
	)
	return salary, err
}
