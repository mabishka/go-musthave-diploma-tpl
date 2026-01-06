package db

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"go.uber.org/zap"
)

func New(ctx context.Context, addr string) (model.Connector, error) {

	logger.Log().Info("address", zap.String("value", addr))
	db, err := sql.Open("pgx", addr)
	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}

	db.SetMaxOpenConns(1000) // Установить максимальное количество открытых соединений к базе данных
	db.SetMaxIdleConns(1000) // Установить максимальное количество неактивных соединений в пуле

	return &model.DB{DB: db}, nil
}
