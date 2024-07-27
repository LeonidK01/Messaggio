package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID uuid.UUID

	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
	DeletedBy *uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	From uuid.UUID
	To   uuid.UUID

	Text   string
	Sended bool
}

// TODO: need add list
type MessageUsecase interface {
	Send(ctx context.Context, msg *Message) error
	ReadByID(ctx context.Context, id uuid.UUID) (*Message, error)
	UpdateByID(ctx context.Context, msg *Message) (*Message, error)
	DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
}

type MessageRepository interface {
	Create(ctx context.Context, msg *Message) (*Message, error)
	ReadByID(ctx context.Context, id uuid.UUID) (*Message, error)
	UpdateByID(ctx context.Context, msg *Message) (*Message, error)
	DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
}

type MessageBroker interface {
	Produce(ctx context.Context, msg *Message) error
}
