package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	service service
	logger  *zap.Logger
}

func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var req CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		h.logger.Error("Failed to decode request", zap.Error(err))
		return
	}

	id, err := h.service.ProcessMessage(r.Context(), req.Content)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		h.logger.Error("Failed to process message", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateMessageResponse{ID: id})
}
