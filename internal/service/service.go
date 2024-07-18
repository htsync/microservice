package service

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/htsync/microservice/tree/main/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	repo          *repository.Repository
	kafkaProducer sarama.AsyncProducer
	logger        *zap.Logger
}

func NewService(repo *repository.Repository, kafkaProducer sarama.AsyncProducer, logger *zap.Logger) *Service {
	return &Service{repo: repo, kafkaProducer: kafkaProducer, logger: logger}
}

type Message struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

type Statistics struct {
	ProcessedMessages int `json:"processed_messages"`
}

func (s *Service) ProcessMessage(ctx context.Context, msg *Message) error {
	s.logger.Info("Processing message", zap.Int("message_id", msg.ID))

	msg.Status = "new"
	if err := s.repo.SaveMessage(ctx, (*repository.Message)(msg)); err != nil {
		return err
	}

	go s.sendToKafka(ctx, msg)
	return nil
}

func (s *Service) sendToKafka(ctx context.Context, msg *Message) {
	select {
	case s.kafkaProducer.Input() <- &sarama.ProducerMessage{
		Topic: "messages",
		Value: sarama.StringEncoder(msg.Content),
	}:
		s.logger.Info("Message sent to Kafka", zap.Int("message_id", msg.ID))
		s.repo.UpdateMessageStatus(ctx, msg.ID, "processed")
	case err := <-s.kafkaProducer.Errors():
		s.logger.Error("Failed to send message to Kafka", zap.Error(err))
		s.repo.UpdateMessageStatus(ctx, msg.ID, "failed")
	}
}

func (s *Service) GetStatistics(ctx context.Context) (*Statistics, error) {
	count, err := s.repo.GetProcessedMessagesCount(ctx)
	if err != nil {
		return nil, err
	}
	return &Statistics{ProcessedMessages: count}, nil
}
