package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"go.uber.org/zap"
)

func (s *Server) HandlerGetBalance(w http.ResponseWriter, r *http.Request) {
	logger.Log().Info("HandlerGetBalance start")
	if r.Method != http.MethodGet {
		logger.Log().Error("HandlerGetBalance error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := r.Header.Get(model.HeaderAuth)
	user, _, err := s.GetUser(r.Context(), token)
	if err != nil {
		if errors.Is(err, model.ErrorNotFound) {
			logger.Log().Error("HandlerGetBalance error GetUser", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		logger.Log().Error("HandlerGetBalance error GetUser", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	current, withdrawn, err := s.GetBalance(r.Context(), user)
	if err != nil {
		logger.Log().Error("HandlerGetBalance error GetBalance", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balance := model.BalanceResponse{
		Current:   current,
		Withdrawn: withdrawn,
	}

	response, err := json.Marshal(balance)
	if err != nil {
		logger.Log().Error("HandlerGetBalance error Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Log().Info("HandlerGetBalance finish ok", zap.Int("status", http.StatusOK))
	w.Header().Set(model.HeaderContentType, model.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
