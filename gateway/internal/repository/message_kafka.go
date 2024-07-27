package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/LeonidK01/Messaggio/internal/model"
	"github.com/segmentio/kafka-go"
)

type messageKafkaBroker struct {
	kWriter *kafka.Writer
}

func NewMessageKafkaBroker(kw *kafka.Writer) model.MessageBroker {
	return &messageKafkaBroker{
		kWriter: kw,
	}
}

func (r *messageKafkaBroker) Produce(ctx context.Context, msg *model.Message) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed marshal message: %w", err)
	}

	kafkaMsg := kafka.Message{
		Key:   []byte(fmt.Sprintf("message-%s", msg.ID.String())),
		Value: bytes,
	}

	if err := r.kWriter.WriteMessages(ctx, kafkaMsg); err != nil {
		return fmt.Errorf("failed write kafka message: %w", err)
	}

	return nil
}
