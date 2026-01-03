package order

import (
	"context"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/repository/db"
	"github.com/mabishka/go-musthave-diploma-tpl/pkg/luhn"
	"go.uber.org/zap"
)

type OrderService interface {
	Add(ctx context.Context, order int, user int) error
	GetOrderList(ctx context.Context, user int) ([]model.OrderResponse, error)
	GetWithdrawals(ctx context.Context, user int) ([]model.WithdrawnResponse, error)
	GetBalance(ctx context.Context, user int) (float32, float32, error)
	GetActiveOrderList(ctx context.Context) ([]model.AccrualProcess, error)
	Withdraw(ctx context.Context, user int, order int, sum float32) error
}

type OrderData struct {
	conn model.Connector
}

func New(ctx context.Context, conn model.Connector) (*OrderData, error) {

	if err := conn.PingContext(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}
	if err := db.CreateOrder(ctx, conn); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err

	}

	return &OrderData{conn: conn}, nil
}

func (p *OrderData) Add(ctx context.Context, order int, user int) error {
	if !luhn.Valid(order) {
		return model.ErrorInvalidLuhn
	}
	if err := p.conn.PingContext(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return err
	}

	logger.Log().Info("ProcessOrder", zap.Int("order", order), zap.Int("user", user))
	if err := db.ProcessOrder(ctx, p.conn, order, user); err != nil {
		logger.Log().Error("error ProcessOrder", zap.Error(err))
		return err
	}

	logger.Log().Info("ProcessOrder ok", zap.Int("order", order), zap.Int("user", user))
	return nil
}

func (p *OrderData) GetOrderList(ctx context.Context, user int) ([]model.OrderResponse, error) {
	if err := p.conn.PingContext(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}
	return db.GetOrderList(ctx, p.conn, user)
}

func (p *OrderData) GetActiveOrderList(ctx context.Context) ([]model.AccrualProcess, error) {
	if err := p.conn.PingContext(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}
	return db.GetActiveOrderList(ctx, p.conn)
}

func (p *OrderData) GetWithdrawals(ctx context.Context, user int) ([]model.WithdrawnResponse, error) {
	if err := p.conn.PingContext(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}
	return db.GetWithdrawals(ctx, p.conn, user)
}

func (p *OrderData) GetBalance(ctx context.Context, user int) (float32, float32, error) {
	if err := p.conn.PingContext(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return 0, 0, err
	}
	return db.GetBalance(ctx, p.conn, user)
}

func (p *OrderData) Withdraw(ctx context.Context, user int, order int, sum float32) error {
	if !luhn.Valid(order) {
		return model.ErrorInvalidLuhn
	}
	if err := p.conn.PingContext(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return err
	}

	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	current, _, err := db.GetBalance(ctx, tx, user)
	if err != nil {
		tx.Rollback()
		return err
	}

	if current < sum {
		tx.Rollback()
		return model.ErrorInvalidBalance
	}

	if err := db.Withdraw(ctx, tx, user, order, sum); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
