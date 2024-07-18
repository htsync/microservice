package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/htsync/microservice/tree/main/internal/handler"
	"github.com/htsync/microservice/tree/main/internal/repository"
	"github.com/htsync/microservice/tree/main/internal/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("Starting the application")

	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal("Failed to connect to the database", zap.Error(err))
	}
	defer db.Close()

	kafkaProducer, err := sarama.NewAsyncProducer([]string{os.Getenv("KAFKA_BROKER")}, nil)
	if err != nil {
		logger.Fatal("Failed to create Kafka producer", zap.Error(err))
	}
	defer kafkaProducer.AsyncClose()

	repo := repository.NewRepository(db, logger)
	svc := service.NewService(repo, kafkaProducer, logger)
	h := handler.NewHandler(svc, logger)

	server := &http.Server{
		Addr:    ":8080",
		Handler: h.InitRoutes(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP server ListenAndServe", zap.Error(err))
		}
	}()

	logger.Info("HTTP server started on :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
