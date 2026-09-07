package models

import (
	"context"
	"database/sql"
	"time"
)

type Expense struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Amount      float64   `json:"amount"`
	ExpenseDate time.Time `json:"expense_date"`
	Notes       *string   `json:"notes,omitempty"`
	ImageURL    *string   `json:"image_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ExpenseModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *ExpenseModel) Insert(ctx context.Context, expense *Expense) (int64, error) {
	query := `
		INSERT INTO expenses (user_id, title, amount, expense_date, notes, image_url)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		expense.UserID,
		expense.Title,
		expense.Amount,
		expense.ExpenseDate,
		expense.Notes,
		expense.ImageURL,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *ExpenseModel) GetByID(ctx context.Context, id int64) (*Expense, error) {
	query := `
		SELECT id, user_id, title, amount, expense_date, notes, image_url, created_at, updated_at
		FROM expenses
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanExpense(row)
}

func (m *ExpenseModel) GetByUserID(ctx context.Context, userID int64) ([]Expense, error) {
	query := `
		SELECT id, user_id, title, amount, expense_date, notes, image_url, created_at, updated_at
		FROM expenses
		WHERE user_id = ?
		ORDER BY expense_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := []Expense{}
	for rows.Next() {
		expense, err := m.scanExpenseRow(rows)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, *expense)
	}

	return expenses, rows.Err()
}

func (m *ExpenseModel) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]Expense, error) {
	query := `
		SELECT id, user_id, title, amount, expense_date, notes, image_url, created_at, updated_at
		FROM expenses
		WHERE expense_date >= ? AND expense_date <= ?
		ORDER BY expense_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := []Expense{}
	for rows.Next() {
		expense, err := m.scanExpenseRow(rows)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, *expense)
	}

	return expenses, rows.Err()
}

func (m *ExpenseModel) GetAll(ctx context.Context) ([]Expense, error) {
	query := `
		SELECT id, user_id, title, amount, expense_date, notes, image_url, created_at, updated_at
		FROM expenses
		ORDER BY expense_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := []Expense{}
	for rows.Next() {
		expense, err := m.scanExpenseRow(rows)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, *expense)
	}

	return expenses, rows.Err()
}

func (m *ExpenseModel) Search(ctx context.Context, query string) ([]Expense, error) {
	searchQuery := `
		SELECT id, user_id, title, amount, expense_date, notes, image_url, created_at, updated_at
		FROM expenses
		WHERE title LIKE ? OR notes LIKE ?
		ORDER BY expense_date DESC
	`

	searchTerm := "%" + query + "%"

	rows, err := m.DB.QueryContext(ctx, searchQuery, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := []Expense{}
	for rows.Next() {
		expense, err := m.scanExpenseRow(rows)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, *expense)
	}

	return expenses, rows.Err()
}

func (m *ExpenseModel) GetTotalByDateRange(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE expense_date >= ? AND expense_date <= ?`

	var total float64
	err := m.DB.QueryRowContext(ctx, query, startDate, endDate).Scan(&total)
	return total, err
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *ExpenseModel) Update(ctx context.Context, expense *Expense) error {
	query := `
		UPDATE expenses
		SET 
			title = ?,
			amount = ?,
			expense_date = ?,
			notes = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		expense.Title,
		expense.Amount,
		expense.ExpenseDate,
		expense.Notes,
		expense.ID,
	)
	return err
}

func (m *ExpenseModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	query := `
		UPDATE expenses
		SET image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, imageURL, id)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *ExpenseModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM expenses WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *ExpenseModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM expenses WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *ExpenseModel) scanExpense(row *sql.Row) (*Expense, error) {
	expense := &Expense{}
	err := row.Scan(
		&expense.ID,
		&expense.UserID,
		&expense.Title,
		&expense.Amount,
		&expense.ExpenseDate,
		&expense.Notes,
		&expense.ImageURL,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return expense, nil
}

func (m *ExpenseModel) scanExpenseRow(rows *sql.Rows) (*Expense, error) {
	expense := &Expense{}
	err := rows.Scan(
		&expense.ID,
		&expense.UserID,
		&expense.Title,
		&expense.Amount,
		&expense.ExpenseDate,
		&expense.Notes,
		&expense.ImageURL,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)
	return expense, err
}
