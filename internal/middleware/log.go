package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"go.uber.org/zap"
)

func WithLogging(h http.Handler) http.Handler {

	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &logger.ResponseData{}
		lw := logger.LoggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			ResponseData:   responseData,
		}
		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		logger.Log().Info("statistic",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
			zap.Int("status", responseData.Status), // получаем перехваченный код статуса ответа
			zap.Int("size", responseData.Size),     // получаем перехваченный размер ответа
			zap.String("headers", fmt.Sprintf("%+v", responseData.Headers)),
		)

	}
	return http.HandlerFunc(logFn)
}
