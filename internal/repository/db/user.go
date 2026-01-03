package db

import (
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"go.uber.org/zap"
)

func CreateUser(ctx context.Context, e model.Executor) error {

	_, err := e.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS t_user (
			n_id SERIAL PRIMARY KEY, 
    		s_login VARCHAR(100) NOT NULL,
    		s_auth VARCHAR(1000) NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_user_login ON t_user(s_login);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_user_auth ON t_user(s_auth);`)

	if err != nil {
		logger.Log().Error("error", zap.Error(err))
	}
	return err
}

func AddUser(ctx context.Context, e model.Executor, login, auth string) error {

	_, err := e.ExecContext(ctx, `insert into t_user(s_login, s_auth) values($1, $2)`, login, auth)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			logger.Log().Error("error", zap.Error(err))
			return model.ErrorAlreadyExist
		}
		logger.Log().Error("error", zap.Error(err))
		return err
	}
	return err
}

func GetAuth(ctx context.Context, e model.Executor, login string) (int, string, error) {

	rows, err := e.QueryContext(ctx,
		`select n_id, s_auth from t_user where s_login = $1`, login)

	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		return 0, "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, "", model.ErrorNotFound
	}

	var auth string
	var id int
	if err := rows.Scan(&id, &auth); err != nil {
		return 0, "", err
	}

	if err = rows.Err(); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return 0, "", err
	}

	return id, auth, nil
}

func GetUser(ctx context.Context, e model.Executor, auth string) (int, string, error) {

	rows, err := e.QueryContext(ctx,
		`select n_id, s_login from t_user where s_auth = $1`, auth)

	if err != nil {
		logger.Log().Error("error GetUser", zap.Error(err), zap.String("auth", auth))
		return 0, "", err
	}
	defer rows.Close()

	if !rows.Next() {
		logger.Log().Error("error GetUser not found", zap.String("auth", auth))
		return 0, "", model.ErrorNotFound
	}

	var login string
	var id int
	if err := rows.Scan(&id, &login); err != nil {
		return 0, "", err
	}

	if err = rows.Err(); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return 0, "", err
	}

	return id, login, nil
}
