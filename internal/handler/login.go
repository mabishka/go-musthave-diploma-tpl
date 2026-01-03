package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/auth"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"go.uber.org/zap"
)

func (s *Server) HandlerPostLogin(w http.ResponseWriter, r *http.Request) {
	logger.Log().Info("HandlerPostLogin start")
	if r.Method != http.MethodPost {
		logger.Log().Error("HandlerPostLogin error method")
		w.WriteHeader(http.StatusBadRequest)
	}

	contentType := r.Header.Get(model.HeaderContentType)
	if contentType != model.ContentTypeJSON {
		logger.Log().Error("HandlerPostLogin error contentType")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log().Error("HandlerPostLogin error ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var register model.RegisterRequest
	if err := json.Unmarshal(body, &register); err != nil {
		logger.Log().Error("HandlerPostLogin error Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := auth.NewToken(register.Login, register.Password)
	if err != nil {
		logger.Log().Error("HandlerPostLogin error NewToken", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	_, login, err := s.GetUser(context.TODO(), token)
	if err != nil {
		if errors.Is(err, model.ErrorNotFound) {
			logger.Log().Error("HandlerPostLogin error GetUser", zap.Error(err), zap.String("login", register.Login))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		logger.Log().Error("HandlerPostLogin error GetUser", zap.Error(err), zap.String("login", register.Login))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if login != register.Login {
		logger.Log().Error("HandlerPostLogin error Login")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	logger.Log().Info("HandlerPostLogin finish ok", zap.Int("status", http.StatusOK), zap.String("token", token))
	w.Header().Set(model.HeaderAuth, token)
	w.WriteHeader(http.StatusOK)
}
