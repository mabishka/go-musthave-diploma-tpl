package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/auth"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"go.uber.org/zap"
)

func (s *Server) HandlerPostRegister(w http.ResponseWriter, r *http.Request) {
	logger.Log().Info("HandlerPostRegister start")
	if r.Method != http.MethodPost {
		logger.Log().Error("HandlerPostRegister error method")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get(model.HeaderContentType)
	if contentType != model.ContentTypeJSON {
		logger.Log().Error("HandlerPostRegister error contentType")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log().Error("HandlerPostRegister error ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var request model.RegisterRequest
	if err := json.Unmarshal(body, &request); err != nil {
		logger.Log().Error("HandlerPostRegister error Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := auth.NewToken(request.Login, request.Password)
	if err != nil {
		logger.Log().Error("HandlerPostRegister error NewToken", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.AddUser(r.Context(), request.Login, token); err != nil {
		if errors.Is(err, model.ErrorAlreadyExist) {

			_, existToken, err := s.GetAuth(r.Context(), request.Login)
			if err != nil {
				logger.Log().Error("HandlerPostRegister error NewToken", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			logger.Log().Error("HandlerPostRegister error AddUser", zap.Error(err))
			if existToken != token {
				w.WriteHeader(http.StatusConflict)
				return
			}
			w.Header().Set(model.HeaderAuth, token)
			w.WriteHeader(http.StatusOK)
			return

		}
		logger.Log().Error("HandlerPostRegister error AddUser", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Log().Info("HandlerPostRegister finish ok", zap.Int("status", http.StatusOK), zap.String("token", token))
	w.Header().Set(model.HeaderAuth, token)
	w.WriteHeader(http.StatusOK)
}
