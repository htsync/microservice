package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Repository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewRepository(db *sqlx.DB, logger *zap.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

type Message struct {
	ID      int    `db:"id"`
	Content string `db:"content"`
	Status  string `db:"status"`
}

func (r *Repository) SaveMessage(ctx context.Context, msg *Message) error {
	_, err := r.db.NamedExecContext(ctx, "INSERT INTO messages (content, status) VALUES (:content, :status)", msg)
	if err != nil {
		r.logger.Error("Failed to save message", zap.Error(err))
	}
	return err
}

func (r *Repository) UpdateMessageStatus(ctx context.Context, id int, status string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE messages SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		r.logger.Error("Failed to update message status", zap.Error(err))
	}
	return err
}

func (r *Repository) GetProcessedMessagesCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM messages WHERE status = 'processed'")
	if err != nil {
		r.logger.Error("Failed to get processed messages count", zap.Error(err))
	}
	return count, err
}
