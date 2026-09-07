package models

import (
	"context"
	"database/sql"
	"time"
)

type Log struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	EntityType  string    `json:"entity_type"`
	EntityID    int64     `json:"entity_id"`
	OldData     *string   `json:"old_data"`
	NewData     *string   `json:"new_data"`
	IPAddress   *string   `json:"ip_address"`
	UserAgent   *string   `json:"user_agent"`
	CreatedAt   time.Time `json:"created_at"`
}

type LogModel struct {
	DB *sql.DB
}

func (m *LogModel) Insert(ctx context.Context, log *Log) (int64, error) {
	query := `
		INSERT INTO logs (user_id, action, description, entity_type, entity_id, old_data, new_data, ip_address, user_agent)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := m.DB.ExecContext(ctx, query,
		log.UserID,
		log.Action,
		log.Description,
		log.EntityType,
		log.EntityID,
		log.OldData,
		log.NewData,
		log.IPAddress,
		log.UserAgent,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *LogModel) GetAll(ctx context.Context) ([]Log, error) {
	query := `
		SELECT id, user_id, action, description, entity_type, entity_id, old_data, new_data, ip_address, user_agent, created_at
		FROM logs
		ORDER BY id DESC
		LIMIT 500
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []Log{}

	for rows.Next() {
		var l Log

		err := rows.Scan(
			&l.ID, &l.UserID, &l.Action, &l.Description,
			&l.EntityType, &l.EntityID, &l.OldData, &l.NewData,
			&l.IPAddress, &l.UserAgent, &l.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		logs = append(logs, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (m *LogModel) GetByEntity(ctx context.Context, entityType string, entityID int64) ([]Log, error) {
	query := `
		SELECT id, user_id, action, description, entity_type, entity_id, old_data, new_data, ip_address, user_agent, created_at
		FROM logs
		WHERE entity_type = ? AND entity_id = ?
		ORDER BY id DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []Log{}

	for rows.Next() {
		var l Log

		err := rows.Scan(
			&l.ID, &l.UserID, &l.Action, &l.Description,
			&l.EntityType, &l.EntityID, &l.OldData, &l.NewData,
			&l.IPAddress, &l.UserAgent, &l.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		logs = append(logs, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
