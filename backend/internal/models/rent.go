package models

import (
	"context"
	"database/sql"
	"time"
)

type Rent struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"user_id"`
	GodownID        int64      `json:"godown_id"`
	MonthYear       string     `json:"month_year"`
	Amount          float64    `json:"amount"`
	Status          string     `json:"status"` // draft, paid, cancelled
	PaymentDate     *time.Time `json:"payment_date,omitempty"`
	PaymentMethod   *string    `json:"payment_method,omitempty"`
	ReferenceNumber *string    `json:"reference_number,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type RentModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *RentModel) Insert(ctx context.Context, rent *Rent) (int64, error) {
	query := `
		INSERT INTO rents (
			user_id, godown_id, month_year, amount, status, notes
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		rent.UserID,
		rent.GodownID,
		rent.MonthYear,
		rent.Amount,
		"draft",
		rent.Notes,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *RentModel) GetByID(ctx context.Context, id int64) (*Rent, error) {
	query := `
		SELECT 
			id, user_id, godown_id, month_year, amount,
			status, payment_date, payment_method, reference_number, notes,
			created_at, updated_at
		FROM rents
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanRent(row)
}

func (m *RentModel) GetByGodownID(ctx context.Context, godownID int64) ([]Rent, error) {
	query := `
		SELECT 
			id, user_id, godown_id, month_year, amount,
			status, payment_date, payment_method, reference_number, notes,
			created_at, updated_at
		FROM rents
		WHERE godown_id = ?
		ORDER BY month_year DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, godownID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rents := []Rent{}
	for rows.Next() {
		rent, err := m.scanRentRow(rows)
		if err != nil {
			return nil, err
		}
		rents = append(rents, *rent)
	}

	return rents, rows.Err()
}

func (m *RentModel) GetByGodownAndMonth(ctx context.Context, godownID int64, monthYear string) (*Rent, error) {
	query := `
		SELECT 
			id, user_id, godown_id, month_year, amount,
			status, payment_date, payment_method, reference_number, notes,
			created_at, updated_at
		FROM rents
		WHERE godown_id = ? AND month_year = ?
	`

	row := m.DB.QueryRowContext(ctx, query, godownID, monthYear)
	return m.scanRent(row)
}

func (m *RentModel) GetByMonth(ctx context.Context, monthYear string) ([]Rent, error) {
	query := `
		SELECT 
			id, user_id, godown_id, month_year, amount,
			status, payment_date, payment_method, reference_number, notes,
			created_at, updated_at
		FROM rents
		WHERE month_year = ?
		ORDER BY godown_id ASC
	`

	rows, err := m.DB.QueryContext(ctx, query, monthYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rents := []Rent{}
	for rows.Next() {
		rent, err := m.scanRentRow(rows)
		if err != nil {
			return nil, err
		}
		rents = append(rents, *rent)
	}

	return rents, rows.Err()
}

func (m *RentModel) GetByStatus(ctx context.Context, status string) ([]Rent, error) {
	query := `
		SELECT 
			id, user_id, godown_id, month_year, amount,
			status, payment_date, payment_method, reference_number, notes,
			created_at, updated_at
		FROM rents
		WHERE status = ?
		ORDER BY month_year DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rents := []Rent{}
	for rows.Next() {
		rent, err := m.scanRentRow(rows)
		if err != nil {
			return nil, err
		}
		rents = append(rents, *rent)
	}

	return rents, rows.Err()
}

func (m *RentModel) GetAll(ctx context.Context) ([]Rent, error) {
	query := `
		SELECT 
			id, user_id, godown_id, month_year, amount,
			status, payment_date, payment_method, reference_number, notes,
			created_at, updated_at
		FROM rents
		ORDER BY month_year DESC, godown_id ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rents := []Rent{}
	for rows.Next() {
		rent, err := m.scanRentRow(rows)
		if err != nil {
			return nil, err
		}
		rents = append(rents, *rent)
	}

	return rents, rows.Err()
}

func (m *RentModel) GetOutstanding(ctx context.Context) ([]Rent, error) {
	query := `
		SELECT 
			id, user_id, godown_id, month_year, amount,
			status, payment_date, payment_method, reference_number, notes,
			created_at, updated_at
		FROM rents
		WHERE status IN ('draft')
		ORDER BY month_year ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rents := []Rent{}
	for rows.Next() {
		rent, err := m.scanRentRow(rows)
		if err != nil {
			return nil, err
		}
		rents = append(rents, *rent)
	}

	return rents, rows.Err()
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *RentModel) Update(ctx context.Context, rent *Rent) error {
	query := `
		UPDATE rents
		SET 
			amount = ?,
			notes = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND status = 'draft'
	`

	result, err := m.DB.ExecContext(ctx, query,
		rent.Amount,
		rent.Notes,
		rent.ID,
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

func (m *RentModel) MarkAsPaid(ctx context.Context, id int64, paymentDate *time.Time, paymentMethod, referenceNumber *string) error {
	query := `
		UPDATE rents
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

func (m *RentModel) Cancel(ctx context.Context, id int64) error {
	query := `
		UPDATE rents
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

func (m *RentModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM rents WHERE id = ? AND status = 'draft'`

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

func (m *RentModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM rents WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *RentModel) ExistsForGodownMonth(ctx context.Context, godownID int64, monthYear string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM rents WHERE godown_id = ? AND month_year = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, godownID, monthYear).Scan(&exists)
	return exists, err
}

func (m *RentModel) EnsureCurrentMonthRents(ctx context.Context, userID int64) error {
	currentMonth := time.Now().Format("2006-01")

	godownModel := GodownModel{DB: m.DB}
	godowns, err := godownModel.GetActive(ctx)
	if err != nil {
		return err
	}

	for _, godown := range godowns {
		existing, err := m.GetByGodownAndMonth(ctx, godown.ID, currentMonth)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if existing == nil && godown.MonthlyRent > 0 {
			rent := &Rent{
				UserID:    userID,
				GodownID:  godown.ID,
				MonthYear: currentMonth,
				Amount:    godown.MonthlyRent,
				Status:    "draft",
				Notes:     nil,
			}

			_, err := m.Insert(ctx, rent)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *RentModel) GetTotalOutstanding(ctx context.Context) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM rents WHERE status = 'draft'`

	var total float64
	err := m.DB.QueryRowContext(ctx, query).Scan(&total)
	return total, err
}

func (m *RentModel) GetTotalByGodown(ctx context.Context, godownID int64) (float64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM rents WHERE godown_id = ? AND status = 'paid'`

	var total float64
	err := m.DB.QueryRowContext(ctx, query, godownID).Scan(&total)
	return total, err
}

// =============================================================================
// SCANNERS
// =============================================================================

func (m *RentModel) scanRent(row *sql.Row) (*Rent, error) {
	rent := &Rent{}
	err := row.Scan(
		&rent.ID,
		&rent.UserID,
		&rent.GodownID,
		&rent.MonthYear,
		&rent.Amount,
		&rent.Status,
		&rent.PaymentDate,
		&rent.PaymentMethod,
		&rent.ReferenceNumber,
		&rent.Notes,
		&rent.CreatedAt,
		&rent.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rent, nil
}

func (m *RentModel) scanRentRow(rows *sql.Rows) (*Rent, error) {
	rent := &Rent{}
	err := rows.Scan(
		&rent.ID,
		&rent.UserID,
		&rent.GodownID,
		&rent.MonthYear,
		&rent.Amount,
		&rent.Status,
		&rent.PaymentDate,
		&rent.PaymentMethod,
		&rent.ReferenceNumber,
		&rent.Notes,
		&rent.CreatedAt,
		&rent.UpdatedAt,
	)
	return rent, err
}
