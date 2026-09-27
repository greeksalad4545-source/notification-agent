package storage

import (
	"context"
	"errors"
	"fmt"

	"notification-agent/models"

	"github.com/jackc/pgx/v5"
)

type Postgres struct {
	conn *pgx.Conn
}

func NewPostgres(ctx context.Context, databaseURL string) (*Postgres, error) {
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to PostgreSQL: %w",
			err,
		)
	}

	return &Postgres{
		conn: conn,
	}, nil
}

func (p *Postgres) Close(ctx context.Context) error {
	return p.conn.Close(ctx)
}

func (p *Postgres) Ping(ctx context.Context) error {
	return p.conn.Ping(ctx)
}

func (p *Postgres) CreateTables(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS notifications (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		channel TEXT NOT NULL,
		subject TEXT,
		message TEXT NOT NULL,
		status TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		processed_at TIMESTAMP
	);
	`

	_, err := p.conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf(
			"failed to create notifications table: %w",
			err,
		)
	}

	return nil
}

func (p *Postgres) SaveNotification(
	ctx context.Context,
	notification models.Notification,
	status string,
) error {
	query := `
	INSERT INTO notifications (
		id,
		user_id,
		channel,
		subject,
		message,
		status,
		processed_at
	)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		CASE
			WHEN $6 = 'processed' THEN CURRENT_TIMESTAMP
			ELSE NULL
		END
	)
	ON CONFLICT (id)
	DO UPDATE SET
		status = EXCLUDED.status,
		processed_at = EXCLUDED.processed_at;
	`

	_, err := p.conn.Exec(
		ctx,
		query,
		notification.ID,
		notification.UserID,
		notification.Channel,
		notification.Subject,
		notification.Message,
		status,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to save notification: %w",
			err,
		)
	}

	return nil
}

func (p *Postgres) GetNotificationStatus(
	ctx context.Context,
	notificationID string,
) (string, string, string, error) {
	query := `
	SELECT status, channel, user_id
	FROM notifications
	WHERE id = $1;
	`

	var status string
	var channel string
	var userID string

	err := p.conn.QueryRow(
		ctx,
		query,
		notificationID,
	).Scan(
		&status,
		&channel,
		&userID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", "", pgx.ErrNoRows
		}

		return "", "", "", fmt.Errorf(
			"failed to get notification status: %w",
			err,
		)
	}

	return status, channel, userID, nil
}
