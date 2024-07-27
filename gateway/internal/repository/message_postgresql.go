package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/LeonidK01/Messaggio/internal/model"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type messagePostrgresqlRepository struct {
	conn   *pgx.Conn
	genSQL squirrel.StatementBuilderType
}

func NewMessagePostgresqlRepository(conn *pgx.Conn) model.MessageRepository {
	return &messagePostrgresqlRepository{
		conn:   conn,
		genSQL: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *messagePostrgresqlRepository) Create(ctx context.Context, msg *model.Message) (*model.Message, error) {
	msg.ID = uuid.New()

	query, args, err := r.genSQL.Insert("message").
		SetMap(map[string]any{
			"id":         msg.ID.String(),
			"created_by": msg.CreatedBy,
			"updated_by": msg.CreatedBy,
			"from":       msg.From.String(),
			"text":       msg.Text,
			"to":         msg.To.String(),
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed generate sql: %w", err)
	}

	if _, err := r.conn.Exec(ctx, query, args...); err != nil {
		return nil, fmt.Errorf("failed exec query: %w", err)
	}

	return msg, nil
}

func (r *messagePostrgresqlRepository) ReadByID(ctx context.Context, id uuid.UUID) (*model.Message, error) {
	query, args, err := r.genSQL.
		Select(
			"id",
			"created_by",
			"updated_by",
			"deleted_by",
			"created_at",
			"updated_at",
			"deleted_at",
			"from",
			"to",
			"text",
			"sended",
		).
		From("message").
		Where(squirrel.And{
			squirrel.Eq{
				"id": id.String(),
			},
			squirrel.Eq{
				"deleted_at": nil,
			},
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed generate sql: %w", err)
	}

	msg := &model.Message{}

	if err := r.conn.QueryRow(ctx, query, args...).Scan(
		&msg.ID,
		&msg.CreatedBy,
		&msg.UpdatedBy,
		&msg.DeletedBy,
		&msg.CreatedAt,
		&msg.UpdatedAt,
		&msg.DeletedAt,
		&msg.From,
		&msg.To,
		&msg.Text,
		&msg.Sended,
	); err != nil {
		return nil, fmt.Errorf("failed query row sql: %w", err)
	}

	return msg, nil
}

func (r *messagePostrgresqlRepository) UpdateByID(ctx context.Context, msg *model.Message) (*model.Message, error) {
	msg.UpdatedAt = time.Now()

	query, args, err := r.genSQL.Update("message").
		Where(squirrel.Eq{
			"id": msg.ID.String(),
		}).
		SetMap(map[string]any{
			"updated_by": msg.UpdatedBy.String(),
			"updated_at": msg.UpdatedAt,
			"text":       msg.Text,
			"sended":     msg.Sended,
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed generate sql: %w", err)
	}

	if _, err := r.conn.Exec(ctx, query, args...); err != nil {
		return nil, fmt.Errorf("failed exec query: %w", err)
	}

	return msg, nil
}

func (r *messagePostrgresqlRepository) DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	query, args, err := r.genSQL.
		Update("message").
		Where(squirrel.Eq{
			"id": id,
		}).
		SetMap(map[string]any{
			"deleted_by": deletedBy.String(),
			"deleted_at": time.Now(),
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed generate sql: %w", err)
	}

	if _, err := r.conn.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("failed exec query: %w", err)
	}

	return nil
}
