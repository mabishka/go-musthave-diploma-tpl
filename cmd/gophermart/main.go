package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/config"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/handler"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/middleware"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/repository/db"
)

const stopTimeout = 5 * time.Second

func main() {

	if err := new(context.WithCancelCause(context.Background())); err != nil {
		log.Fatalf("exist with error: %v", err)
	}
}

func new(ctx context.Context, fnCancel context.CancelCauseFunc) error {

	config := config.New()
	if err := logger.InitLogger(config.LogLevel); err != nil {
		fnCancel(err)
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
		fnCancel(err)
		return err
	}
	defer conn.Close()

	server, err := handler.New(ctx, conn, config.AccuralSystemAddress)
	if err != nil {
		logger.Log().Error("incorrect server", zap.Error(err))
		fnCancel(err)
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

	if err = run(ctx, &http.Server{
		Addr:         config.RunAddress,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}); err != nil {
		fnCancel(err)
		return err
	}
	fnCancel(nil)

	return nil
}
func run(ctx context.Context, srv *http.Server) error {

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		select {
		case s := <-sigint:
			logger.Log().Info("stop with signal", zap.String("signal", s.String()))
		case <-ctx.Done():
			logger.Log().Info("stop with context", zap.Error(context.Cause(ctx)))
		}

		stopCtx, cancel := context.WithTimeoutCause(context.Background(), stopTimeout, fmt.Errorf("server Shutdown with timeout %v", stopTimeout))
		defer cancel()
		if err := srv.Shutdown(stopCtx); err != nil {
			logger.Log().Info("HTTP server shutdown", zap.Error(err))
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Log().Info("HTTP server ListenAndServe", zap.Error(err))
		return err
	}

	logger.Log().Info("exit")
	return nil
}
