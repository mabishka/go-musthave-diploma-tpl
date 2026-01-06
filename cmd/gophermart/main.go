package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/config"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/handler"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/middleware"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/repository/db"
)

func main() {
	ctx, fnCancel := context.WithCancelCause(context.Background())
	defer fnCancel(errors.New("exist"))

	if err := new(ctx); err != nil {
		log.Printf("exist with error: %v", err)
	}
}

func new(ctx context.Context) error {

	config := config.New()
	if err := logger.InitLogger(config.LogLevel); err != nil {
		return err
	}

	logger.Log().Info("config",
		zap.String("RunAddress", config.RunAddress),
		zap.String("DatabaseURL", config.DatabaseURL),
		zap.String("AccuralServerAddress", config.AccuralSystemAddress),
		zap.String("LogLevel", config.LogLevel),
	)

	config.Load()

	logger.Log().Info("config",
		zap.String("RunAddress", config.RunAddress),
		zap.String("DatabaseURL", config.DatabaseURL),
		zap.String("AccuralServerAddress", config.AccuralSystemAddress),
		zap.String("LogLevel", config.LogLevel),
	)

	conn, err := db.New(ctx, config.DatabaseURL)
	if err != nil {
		logger.Log().Error("incorrect config", zap.Error(err))
		return err
	}
	defer conn.Close()

	server, err := handler.New(ctx, conn, config.AccuralSystemAddress)
	if err != nil {
		logger.Log().Error("incorrect server", zap.Error(err))
		return err
	}
	defer server.Close()

	router := chi.NewRouter()

	router.Use(middleware.WithLogging)
	router.Use(middleware.WithCompress)

	router.Post("/api/user/register", server.HandlerPostRegister)
	router.Post("/api/user/login", server.HandlerPostLogin)
	router.Post("/api/user/orders", server.HandlerPostOrders)
	router.Post("/api/user/balance/withdraw", server.HandlerPostWithdraw)
	router.Get("/api/user/orders", server.HandlerGetOrders)
	router.Get("/api/user/balance", server.HandlerGetBalance)
	router.Get("/api/user/withdrawals", server.HandlerGetWithdrawals)

	run(ctx, &http.Server{
		Addr:    config.RunAddress,
		Handler: router,
	})
	return nil
}

func run(ctx context.Context, srv *http.Server) {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint)

		select {
		case s := <-sigint:
			logger.Log().Info("stop with signal", zap.String("signal", s.String()))
		case <-ctx.Done():
			logger.Log().Info("stop with context", zap.Error(context.Cause(ctx)))
		}

		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Log().Info("HTTP server shutdown", zap.Error(err))
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Log().Info("HTTP server ListenAndServe", zap.Error(err))
	}
	wg.Wait()

	logger.Log().Info("exit")
}
