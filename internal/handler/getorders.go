package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"go.uber.org/zap"
)

func (s *Server) HandlerGetOrders(w http.ResponseWriter, r *http.Request) {
	logger.Log().Info("HandlerGetOrders start")
	if r.Method != http.MethodGet {
		logger.Log().Error("HandlerGetOrders error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := r.Header.Get(model.HeaderAuth)
	user, _, err := s.GetUser(r.Context(), token)
	if err != nil {
		if errors.Is(err, model.ErrorNotFound) {
			logger.Log().Error("HandlerGetOrders error GetUser", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		logger.Log().Error("HandlerGetOrders error GetUser", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	list, err := s.GetOrderList(r.Context(), user)
	if err != nil {
		if errors.Is(err, model.ErrorNotFound) {
			logger.Log().Info("HandlerGetOrders error GetOrderList", zap.Error(err))
			w.WriteHeader(http.StatusNoContent)
			return
		}
		logger.Log().Error("HandlerGetOrders error GetOrderList", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(list) == 0 {
		logger.Log().Info("HandlerGetOrders error GetOrderList is empty")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	logger.Log().Info("HandlerGetOrders get", zap.Int("user", user), zap.Int("len", len(list)))

	response, err := json.Marshal(list)
	if err != nil {
		logger.Log().Error("HandlerGetOrders error Marshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Log().Info("HandlerGetOrders finish ok", zap.Int("status", http.StatusOK))
	w.Header().Set(model.HeaderContentType, model.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
