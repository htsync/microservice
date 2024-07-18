package service_test

import (
	"Microservise/internal/repository"
	"Microservise/internal/service"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestProcessMessage(t *testing.T) {
	logger := zaptest.NewLogger(t)
	repo := repository.NewRepositoryMock()
	kafkaProducer := sarama.NewSyncProducerMock()
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	svc := service.NewService(repo, kafkaProducer, redisClient, logger)

	msg := "test message"
	id, err := svc.ProcessMessage(msg)
	assert.NoError(t, err)
	assert.NotZero(t, id)
}
