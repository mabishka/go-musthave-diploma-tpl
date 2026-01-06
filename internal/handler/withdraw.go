package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"go.uber.org/zap"
)

func (s *Server) HandlerPostWithdraw(w http.ResponseWriter, r *http.Request) {

	logger.Log().Info("HandlerPostWithdraw start")

	if r.Method != http.MethodPost {
		logger.Log().Error("HandlerPostWithdraw error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get(model.HeaderContentType)
	if contentType != model.ContentTypeJSON {
		logger.Log().Error("HandlerPostWithdraw error contentType")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := r.Header.Get(model.HeaderAuth)
	user, _, err := s.GetUser(r.Context(), token)
	if err != nil {
		if errors.Is(err, model.ErrorNotFound) {
			logger.Log().Error("HandlerPostWithdraw error GetUser", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		logger.Log().Error("HandlerPostWithdraw error GetUser", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Log().Info("HandlerPostWithdraw user", zap.Int("user", user))

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log().Error("HandlerPostWithdraw error ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var request model.WithdrawnRequest
	if err := json.Unmarshal(body, &request); err != nil {
		logger.Log().Error("HandlerPostWithdraw error Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Log().Info("HandlerPostWithdraw request", zap.String("order", request.Order))

	order, err := strconv.Atoi(request.Order)
	if err != nil {
		logger.Log().Error("HandlerPostWithdraw error Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.Withdraw(r.Context(), user, order, request.Sum); err != nil {
		if errors.Is(err, model.ErrorInvalidBalance) {
			logger.Log().Error("HandlerPostWithdraw error Withdraw", zap.Error(err))
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}
		if errors.Is(err, model.ErrorInvalidLuhn) || errors.Is(err, model.ErrorNotFound) {
			logger.Log().Error("HandlerPostWithdraw error Withdraw", zap.Error(err))
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		logger.Log().Error("HandlerPostWithdraw error Withdraw", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Log().Info("HandlerPostWithdraw finish ok", zap.Int("status", http.StatusOK))
	w.WriteHeader(http.StatusOK)
}
