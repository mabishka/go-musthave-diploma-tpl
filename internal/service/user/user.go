package user

import (
	"context"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/repository/db"
	"go.uber.org/zap"
)

type UserService interface {
	AddUser(ctx context.Context, login string, token string) error
	GetUser(ctx context.Context, token string) (int, string, error)
	GetAuth(ctx context.Context, login string) (int, string, error)
}

type UserData struct {
	conn model.Connector
}

func New(ctx context.Context, conn model.Connector) (*UserData, error) {
	if err := conn.PingContext(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}
	if err := db.CreateUser(ctx, conn); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err

	}
	return &UserData{conn: conn}, nil
}

func (p *UserData) AddUser(ctx context.Context, login, token string) error {
	if err := p.conn.PingContext(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return err
	}
	return db.AddUser(ctx, p.conn, login, token)
}

func (p *UserData) GetUser(ctx context.Context, token string) (int, string, error) {

	if err := p.conn.PingContext(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return 0, "", err
	}
	return db.GetUser(ctx, p.conn, token)

}

func (p *UserData) GetAuth(ctx context.Context, login string) (int, string, error) {

	if err := p.conn.PingContext(ctx); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return 0, "", err
	}
	return db.GetAuth(ctx, p.conn, login)

}
