package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"go.uber.org/zap"
)

func (s *Server) HandlerGetWithdrawals(w http.ResponseWriter, r *http.Request) {
	logger.Log().Info("HandlerGetWithdrawals start")
	if r.Method != http.MethodGet {
		logger.Log().Error("HandlerGetWithdrawals error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := r.Header.Get(model.HeaderAuth)
	user, _, err := s.GetUser(r.Context(), token)
	if err != nil {
		if errors.Is(err, model.ErrorNotFound) {
			logger.Log().Error("HandlerGetWithdrawals error GetUser", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		logger.Log().Error("HandlerGetWithdrawals error GetUser", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	logger.Log().Info("get withdrawls for user", zap.Int("user", user))
	list, err := s.GetWithdrawals(r.Context(), user)
	if err != nil {
		logger.Log().Error("HandlerGetWithdrawals error GetWithdrawals", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(list) == 0 {
		logger.Log().Error("HandlerGetWithdrawals error len")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response, err := json.Marshal(list)
	if err != nil {
		logger.Log().Error("HandlerGetWithdrawals error Marshal", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Log().Info("HandlerGetWithdrawals finish ok", zap.Int("status", http.StatusOK))
	w.Header().Set(model.HeaderContentType, model.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}
