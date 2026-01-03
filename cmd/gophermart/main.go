package main

import (
	"context"
	"errors"
	"net/http"

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
	defer fnCancel(errors.New("exit"))
	run(ctx)
}

func run(ctx context.Context) {

	config := config.New()
	if err := logger.InitLogger(config.LogLevel); err != nil {
		panic(err)
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
		panic(err)
	}
	defer conn.Close()

	server, err := handler.New(ctx, conn, config.AccuralSystemAddress)
	if err != nil {
		panic(err)
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

	go func() {
		if err := http.ListenAndServe(config.RunAddress, router); err != nil {
			panic(err)
		}
	}()

	logger.Log().Info("listen port", zap.String("RunAddress", config.RunAddress))

	<-ctx.Done()
	logger.Log().Info("exit")
}
