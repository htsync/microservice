package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	// Инициализация логгера
	logger := initLogger()
	defer logger.Sync()

	logger.Info("Starting application...")

	// Создание зависимостей
	dbConnStr := os.Getenv("DB_CONN_STRING")
	repo, err := initRepository(dbConnStr, logger)
	if err != nil {
		logger.Fatal("Failed to initialize repository", zap.Error(err))
	}
	logger.Info("Repository initialized")

	kafkaBrokers := []string{os.Getenv("KAFKA_BROKERS")}
	kafkaProducer, err := initKafkaProducer(kafkaBrokers, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Kafka producer", zap.Error(err))
	}
	defer kafkaProducer.Close()
	logger.Info("Kafka producer initialized")

	redisAddr := os.Getenv("REDIS_ADDR")
	redisClient := initRedisClient(redisAddr, logger)
	logger.Info("Redis client initialized")

	// Создание сервиса с внедрением зависимостей
	svc := initService(repo, kafkaProducer, redisClient, logger)
	logger.Info("Service initialized")

	// Создание HTTP обработчиков с внедрением сервиса
	r := mux.NewRouter()
	h := initHandler(svc, logger)

	r.HandleFunc("/message", h.CreateMessage).Methods("POST")
	r.HandleFunc("/stats", h.GetStats).Methods("GET")

	// Настройка маршрута для метрик Prometheus
	r.Handle("/metrics", promhttp.Handler())

	// Запуск HTTP сервера
	httpPort := 8080
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Ожидание сигнала завершения работы приложения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Info("Shutting down server...")

	// Остановка HTTP сервера с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Failed to gracefully shutdown HTTP server", zap.Error(err))
	}

	logger.Info("Server stopped")
}

func initLogger() *zap.Logger {
	// Инициализация логгера
	config := zap.NewProductionConfig()
	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	return logger
}

func initRepository(dbConnStr string, logger *zap.Logger) (repository, error) {
	// Инициализация репозитория
}

func initKafkaProducer(brokers []string, logger *zap.Logger) (sarama.SyncProducer, error) {
	// Инициализация Kafka producer
}

func initRedisClient(redisAddr string, logger *zap.Logger) *redis.Client {
	// Инициализация клиента Redis
}

func initService(repo repository, kafkaProducer sarama.SyncProducer, redisClient *redis.Client, logger *zap.Logger) service {
	// Инициализация сервиса с внедрением зависимостей
}

func initHandler(svc service, logger *zap.Logger) handler {
	// Инициализация HTTP обработчика с внедрением сервиса
}
