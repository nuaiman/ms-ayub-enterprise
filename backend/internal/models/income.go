package models

import (
	"context"
	"database/sql"
	"time"
)

type Income struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Title      string    `json:"title"`
	Amount     float64   `json:"amount"`
	IncomeDate time.Time `json:"income_date"`
	Notes      *string   `json:"notes,omitempty"`
	ImageURL   *string   `json:"image_url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type IncomeModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *IncomeModel) Insert(ctx context.Context, income *Income) (int64, error) {
	query := `
		INSERT INTO incomes (user_id, title, amount, income_date, notes, image_url)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		income.UserID,
		income.Title,
		income.Amount,
		income.IncomeDate,
		income.Notes,
		income.ImageURL,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *IncomeModel) GetByID(ctx context.Context, id int64) (*Income, error) {
	query := `
		SELECT id, user_id, title, amount, income_date, notes, image_url, created_at, updated_at
		FROM incomes
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanIncome(row)
}

func (m *IncomeModel) GetByUserID(ctx context.Context, userID int64) ([]Income, error) {
	query := `
		SELECT id, user_id, title, amount, income_date, notes, image_url, created_at, updated_at
		FROM incomes
		WHERE user_id = ?
		ORDER BY income_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	incomes := []Income{}
	for rows.Next() {
		income, err := m.scanIncomeRow(rows)
		if err != nil {
			return nil, err
		}
		incomes = append(incomes, *income)
	}

	return incomes, rows.Err()
}

func (m *IncomeModel) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]Income, error) {
	query := `
		SELECT id, user_id, title, amount, income_date, notes, image_url, created_at, updated_at
		FROM incomes
		WHERE income_date >= ? AND income_date <= ?
		ORDER BY income_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	incomes := []Income{}
	for rows.Next() {
		income, err := m.scanIncomeRow(rows)
		if err != nil {
			return nil, err
		}
		incomes = append(incomes, *income)
	}

	return incomes, rows.Err()
}

func (m *IncomeModel) GetAll(ctx context.Context) ([]Income, error) {
	query := `
		SELECT id, user_id, title, amount, income_date, notes, image_url, created_at, updated_at
		FROM incomes
		ORDER BY income_date DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	incomes := []Income{}
	for rows.Next() {
		income, err := m.scanIncomeRow(rows)
		if err != nil {
			return nil, err
		}
		incomes = append(incomes, *income)
	}

	return incomes, rows.Err()
}

func (m *IncomeModel) Search(ctx context.Context, query string) ([]Income, error) {
	searchQuery := `
		SELECT id, user_id, title, amount, income_date, notes, image_url, created_at, updated_at
		FROM incomes
		WHERE title LIKE ? OR notes LIKE ?
		ORDER BY income_date DESC
	`

	searchTerm := "%" + query + "%"

	rows, err := m.DB.QueryContext(ctx, searchQuery, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	incomes := []Income{}
	for rows.Next() {
		income, err := m.scanIncomeRow(rows)
		if err != nil {
			return nil, err
		}
		incomes = append(incomes, *income)
	}

	return incomes, rows.Err()
}

func (m *IncomeModel) GetTotalByDateRange(ctx context.Context, startDate, endDate time.Time) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM incomes WHERE income_date >= ? AND income_date <= ?`

	var total float64
	err := m.DB.QueryRowContext(ctx, query, startDate, endDate).Scan(&total)
	return total, err
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *IncomeModel) Update(ctx context.Context, income *Income) error {
	query := `
		UPDATE incomes
		SET 
			title = ?,
			amount = ?,
			income_date = ?,
			notes = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		income.Title,
		income.Amount,
		income.IncomeDate,
		income.Notes,
		income.ID,
	)
	return err
}

func (m *IncomeModel) UpdateImage(ctx context.Context, id int64, imageURL *string) error {
	query := `
		UPDATE incomes
		SET image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, imageURL, id)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *IncomeModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM incomes WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *IncomeModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM incomes WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *IncomeModel) scanIncome(row *sql.Row) (*Income, error) {
	income := &Income{}
	err := row.Scan(
		&income.ID,
		&income.UserID,
		&income.Title,
		&income.Amount,
		&income.IncomeDate,
		&income.Notes,
		&income.ImageURL,
		&income.CreatedAt,
		&income.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return income, nil
}

func (m *IncomeModel) scanIncomeRow(rows *sql.Rows) (*Income, error) {
	income := &Income{}
	err := rows.Scan(
		&income.ID,
		&income.UserID,
		&income.Title,
		&income.Amount,
		&income.IncomeDate,
		&income.Notes,
		&income.ImageURL,
		&income.CreatedAt,
		&income.UpdatedAt,
	)
	return income, err
}
