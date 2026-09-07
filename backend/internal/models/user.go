package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID                 int64      `json:"id"`
	Name               string     `json:"name"`
	Username           string     `json:"username"`
	Email              *string    `json:"email,omitempty"`
	Phone              *string    `json:"phone,omitempty"`
	Address            *string    `json:"address,omitempty"`
	IDType             *string    `json:"id_type,omitempty"`
	IDNumber           *string    `json:"id_number,omitempty"`
	ImageURL           *string    `json:"image_url,omitempty"`
	Password           string     `json:"-"`
	Role               string     `json:"role"`
	RefreshToken       *string    `json:"-"`
	RefreshTokenExpiry *time.Time `json:"-"`
	IsActive           bool       `json:"is_active"`
	MonthlySalary      float64    `json:"monthly_salary"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type UserModel struct {
	DB *sql.DB
}

// =============================================================================
// CREATE
// =============================================================================

func (m *UserModel) Insert(ctx context.Context, user *User) (int64, error) {
	query := `
		INSERT INTO users (
			name, username, email, phone, address, 
			id_type, id_number, image_url, password, role, monthly_salary
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		user.Name,
		user.Username,
		user.Email,
		user.Phone,
		user.Address,
		user.IDType,
		user.IDNumber,
		user.ImageURL,
		user.Password,
		user.Role,
		user.MonthlySalary,
	)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// =============================================================================
// READ
// =============================================================================

func (m *UserModel) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, name, username, email, phone, address, 
		       id_type, id_number, image_url, password, role, 
		       refresh_token, refresh_token_expiry, is_active, monthly_salary, 
		       created_at, updated_at
		FROM users
		WHERE id = ?
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	return m.scanUser(row)
}

func (m *UserModel) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT id, name, username, email, phone, address, 
		       id_type, id_number, image_url, password, role, 
		       refresh_token, refresh_token_expiry, is_active, monthly_salary, 
		       created_at, updated_at
		FROM users
		WHERE username = ?
	`

	row := m.DB.QueryRowContext(ctx, query, username)
	return m.scanUser(row)
}

func (m *UserModel) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, name, username, email, phone, address, 
		       id_type, id_number, image_url, password, role, 
		       refresh_token, refresh_token_expiry, is_active, monthly_salary, 
		       created_at, updated_at
		FROM users
		WHERE email = ?
	`

	row := m.DB.QueryRowContext(ctx, query, email)
	return m.scanUser(row)
}

func (m *UserModel) GetByRefreshToken(ctx context.Context, token string) (*User, error) {
	query := `
		SELECT id, name, username, email, phone, address, 
		       id_type, id_number, image_url, password, role, 
		       refresh_token, refresh_token_expiry, is_active, monthly_salary, 
		       created_at, updated_at
		FROM users
		WHERE refresh_token = ?
	`

	row := m.DB.QueryRowContext(ctx, query, token)
	user, err := m.scanUser(row)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	// Check if refresh token has expired
	if user.RefreshTokenExpiry != nil && user.RefreshTokenExpiry.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	return user, nil
}

func (m *UserModel) GetAll(ctx context.Context) ([]User, error) {
	query := `
		SELECT id, name, username, email, phone, address, 
		       id_type, id_number, image_url, password, role, 
		       refresh_token, refresh_token_expiry, is_active, monthly_salary, 
		       created_at, updated_at
		FROM users
		ORDER BY id DESC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user, err := m.scanUserRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, rows.Err()
}

func (m *UserModel) scanUser(row *sql.Row) (*User, error) {
	user := &User{}
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.Address,
		&user.IDType,
		&user.IDNumber,
		&user.ImageURL,
		&user.Password,
		&user.Role,
		&user.RefreshToken,
		&user.RefreshTokenExpiry,
		&user.IsActive,
		&user.MonthlySalary,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (m *UserModel) scanUserRow(rows *sql.Rows) (*User, error) {
	user := &User{}
	err := rows.Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.Address,
		&user.IDType,
		&user.IDNumber,
		&user.ImageURL,
		&user.Password,
		&user.Role,
		&user.RefreshToken,
		&user.RefreshTokenExpiry,
		&user.IsActive,
		&user.MonthlySalary,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

// =============================================================================
// UPDATE
// =============================================================================

func (m *UserModel) UpdateRefreshToken(ctx context.Context, userID int64, token *string) error {
	query := `
		UPDATE users
		SET refresh_token = ?, 
		    refresh_token_expiry = CASE WHEN ? IS NOT NULL THEN datetime('now', '+30 days') ELSE NULL END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, token, token, userID)
	return err
}

func (m *UserModel) UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = ?, refresh_token = NULL, refresh_token_expiry = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, hashedPassword, userID)
	return err
}

func (m *UserModel) UpdateProfile(ctx context.Context, userID int64, user *User) error {
	query := `
		UPDATE users
		SET name = ?, email = ?, phone = ?, address = ?, 
		    id_type = ?, id_number = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.Phone,
		user.Address,
		user.IDType,
		user.IDNumber,
		userID,
	)
	return err
}

func (m *UserModel) UpdateImage(ctx context.Context, userID int64, imageURL *string) error {
	query := `
		UPDATE users
		SET image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, imageURL, userID)
	return err
}

func (m *UserModel) UpdateMonthlySalary(ctx context.Context, userID int64, salary float64) error {
	query := `
		UPDATE users
		SET monthly_salary = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, salary, userID)
	return err
}

func (m *UserModel) UpdateRole(ctx context.Context, userID int64, role string) error {
	query := `
		UPDATE users
		SET role = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, role, userID)
	return err
}

func (m *UserModel) SetActiveStatus(ctx context.Context, userID int64, active bool) error {
	isActive := 0
	if active {
		isActive = 1
	}

	query := `
		UPDATE users
		SET is_active = ?, refresh_token = NULL, refresh_token_expiry = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, isActive, userID)
	return err
}

func (m *UserModel) AdminUpdatePasswordAndInvalidateSessions(ctx context.Context, userID int64, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = ?, refresh_token = NULL, refresh_token_expiry = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, hashedPassword, userID)
	return err
}

func (m *UserModel) ResetAllPasswordsExceptRoot(ctx context.Context, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = ?, refresh_token = NULL, refresh_token_expiry = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id != 1
	`

	_, err := m.DB.ExecContext(ctx, query, hashedPassword)
	return err
}

// =============================================================================
// DELETE
// =============================================================================

func (m *UserModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// =============================================================================
// UTILITY
// =============================================================================

func (m *UserModel) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

func (m *UserModel) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, username).Scan(&exists)
	return exists, err
}

func (m *UserModel) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}

func (m *UserModel) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int
	err := m.DB.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (m *UserModel) GetActiveUsers(ctx context.Context) ([]User, error) {
	query := `
		SELECT id, name, username, email, phone, address, 
		       id_type, id_number, image_url, password, role, 
		       refresh_token, refresh_token_expiry, is_active, monthly_salary, 
		       created_at, updated_at
		FROM users
		WHERE is_active = 1
		ORDER BY name ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user, err := m.scanUserRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, rows.Err()
}

func (m *UserModel) GetByRole(ctx context.Context, role string) ([]User, error) {
	query := `
		SELECT id, name, username, email, phone, address, 
		       id_type, id_number, image_url, password, role, 
		       refresh_token, refresh_token_expiry, is_active, monthly_salary, 
		       created_at, updated_at
		FROM users
		WHERE role = ?
		ORDER BY name ASC
	`

	rows, err := m.DB.QueryContext(ctx, query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user, err := m.scanUserRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, rows.Err()
}

func (m *UserModel) Search(ctx context.Context, query string) ([]User, error) {
	searchQuery := `
		SELECT id, name, username, email, phone, address, 
		       id_type, id_number, image_url, password, role, 
		       refresh_token, refresh_token_expiry, is_active, monthly_salary, 
		       created_at, updated_at
		FROM users
		WHERE name LIKE ? OR username LIKE ? OR email LIKE ? OR phone LIKE ?
		ORDER BY name ASC
	`

	searchTerm := "%" + query + "%"

	rows, err := m.DB.QueryContext(ctx, searchQuery,
		searchTerm, searchTerm, searchTerm, searchTerm,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user, err := m.scanUserRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, rows.Err()
}
