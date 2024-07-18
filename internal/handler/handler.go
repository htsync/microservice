package handler

import (
	"encoding/json"
	"net/http"

	"github.com/htsync/microservice/tree/main/internal/service"
	"go.uber.org/zap"
)

type Handler struct {
	service *service.Service
	logger  *zap.Logger
}

func NewHandler(service *service.Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/messages", h.handleMessages)
	mux.HandleFunc("/statistics", h.handleStatistics)
	return mux
}

func (h *Handler) handleMessages(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling messages")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var msg service.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		h.logger.Error("Failed to decode message", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := h.service.ProcessMessage(r.Context(), &msg); err != nil {
		h.logger.Error("Failed to process message", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleStatistics(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handling statistics")

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.service.GetStatistics(r.Context())
	if err != nil {
		h.logger.Error("Failed to get statistics", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		h.logger.Error("Failed to encode statistics", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
