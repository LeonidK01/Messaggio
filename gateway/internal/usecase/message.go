package usecase

import (
	"context"
	"fmt"

	"github.com/LeonidK01/Messaggio/internal/model"
	"github.com/google/uuid"
)

type messageUsecase struct {
	messageRepo   model.MessageRepository
	messageBroker model.MessageBroker
}

func NewMessageUsecase(mr model.MessageRepository, mq model.MessageBroker) model.MessageUsecase {
	return &messageUsecase{
		messageRepo:   mr,
		messageBroker: mq,
	}
}

func (uc *messageUsecase) Send(ctx context.Context, msg *model.Message) error {
	result, err := uc.messageRepo.Create(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed create message in db: %w", err)
	}

	if err := uc.messageBroker.Produce(ctx, result); err != nil {
		// added outbox
		return fmt.Errorf("failed send message in mq: %w", err)
	}

	return nil
}

func (uc *messageUsecase) ReadByID(ctx context.Context, id uuid.UUID) (*model.Message, error) {
	msg, err := uc.messageRepo.ReadByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed read by id: %w", err)
	}

	return msg, nil
}

func (uc *messageUsecase) UpdateByID(ctx context.Context, msg *model.Message) (*model.Message, error) {
	msg, err := uc.messageRepo.UpdateByID(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed update by id: %w", err)
	}

	return msg, nil
}

func (uc *messageUsecase) DeleteByID(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	err := uc.messageRepo.DeleteByID(ctx, id, deletedBy)
	if err != nil {
		return fmt.Errorf("failed delete by id: %w", err)
	}

	return nil
}
