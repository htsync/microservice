package service

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"

	"go.uber.org/zap"

	"github.com/your/module/repository"
)

type Service struct {
	repo          repository.Repository
	kafkaProducer sarama.SyncProducer
	redisClient   *redis.Client
	logger        *zap.Logger
}

func (s *Service) ProcessMessage(ctx context.Context, content string) (int, error) {
	id, err := s.repo.SaveMessage(content)
	if err != nil {
		return 0, err
	}

	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)

		msg := &sarama.ProducerMessage{
			Topic: "messages",
			Value: sarama.StringEncoder(content),
		}

		_, _, err := s.kafkaProducer.SendMessage(msg)
		if err != nil {
			errCh <- err
			return
		}

		if err := s.repo.MarkMessageAsProcessed(id); err != nil {
			errCh <- err
			return
		}

		s.logger.Info("Message processed", zap.Int("id", id))
	}()

	select {
	case <-ctx.Done():
		return 0, errors.New("operation canceled")
	case err := <-errCh:
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}
