package handler

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"go.uber.org/zap"
)

func (s *Server) HandlerPostOrders(w http.ResponseWriter, r *http.Request) {
	logger.Log().Info("HandlerPostOrders start")
	if r.Method != http.MethodPost {
		logger.Log().Error("HandlerPostOrders error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get(model.HeaderContentType)
	if contentType != model.ContentTypeText {
		logger.Log().Error("HandlerPostOrders error contentType")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := r.Header.Get(model.HeaderAuth)
	user, _, err := s.GetUser(r.Context(), token)
	if err != nil {
		if errors.Is(err, model.ErrorNotFound) {
			logger.Log().Error("HandlerPostOrders error GetUser", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		logger.Log().Error("HandlerPostOrders error GetUser", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log().Error("HandlerPostOrders error ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	num, err := strconv.Atoi(string(body))
	if err != nil {
		logger.Log().Error("HandlerPostOrders error Atoi", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.Add(r.Context(), num, user); err != nil {
		if errors.Is(err, model.ErrorInvalidLuhn) {
			logger.Log().Error("HandlerPostOrders error Add", zap.Error(err))
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, model.ErrorAlreadyExist) {
			logger.Log().Info("HandlerPostOrders finish ok", zap.Int("status", http.StatusOK))
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, model.ErrorNotOwned) {
			logger.Log().Error("HandlerPostOrders error Add", zap.Error(err))
			w.WriteHeader(http.StatusConflict)
			return
		}
		logger.Log().Error("HandlerPostOrders error Add", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.ProcessOrder(model.AccrualProcess{User: user, Order: num})

	logger.Log().Info("HandlerPostOrders finish ok", zap.Int("status", http.StatusAccepted))
	w.WriteHeader(http.StatusAccepted)
}
