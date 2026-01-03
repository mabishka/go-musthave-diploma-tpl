package db

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"go.uber.org/zap"
)

func CreateOrder(ctx context.Context, e model.Executor) error {

	_, err := e.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS t_order (
    		n_id INT8 NOT NULL PRIMARY KEY,
			n_user INT NOT NULL,
    		s_status VARCHAR(20),
			s_state VARCHAR(20),
			dt_uploaded timestamptz,
			CONSTRAINT fk_order_user FOREIGN KEY (n_user) REFERENCES t_user(n_id)
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_order_user_id ON t_order(n_id, n_user);
		CREATE TABLE IF NOT EXISTS t_order_value (
			n_user INT NOT NULL,
    		n_order INT8 NOT NULL,
    		n_value NUMERIC(10,2) NOT NULL,
			dt_processed timestamptz
		);`)
	if err != nil {
		logger.Log().Error("error", zap.Error(err))
	}
	return err
}

func ProcessOrder(ctx context.Context, e model.Executor, order, user int) error {

	if _, err := e.ExecContext(ctx, `insert into t_order(n_id, n_user, s_status, s_state, dt_uploaded) values($1, $2, $3, $4, current_timestamp)`,
		order, user, model.AccrualStatusEmpty, model.GetState(model.AccrualStatusEmpty)); err != nil {

		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {

			logger.Log().Error("order UniqueViolation", zap.Error(err))

			rows, err := e.QueryContext(ctx, "select n_user = $1 from t_order where n_id = $2", user, order)
			if err != nil {
				logger.Log().Error("order UniqueViolation select error", zap.Error(err))
				return err
			}
			if !rows.Next() {
				logger.Log().Error("order UniqueViolation select ErrorNotFound")
				return model.ErrorNotFound
			}
			var isOwned bool
			if err := rows.Scan(&isOwned); err != nil {
				logger.Log().Error("order UniqueViolation scan error", zap.Error(err))
				return model.ErrorNotFound
			}
			if err = rows.Err(); err != nil {
				logger.Log().Error("error", zap.Error(err))
				return err
			}
			if !isOwned {
				logger.Log().Error("order UniqueViolation ErrorNotOwned")
				return model.ErrorNotOwned
			}
			logger.Log().Error("order UniqueViolation ErrorAlreadyExist")
			return model.ErrorAlreadyExist
		}

		logger.Log().Error("order other error", zap.Error(err))
		return err
	}
	return nil
}

func GetBalance(ctx context.Context, e model.Executor, user int) (float32, float32, error) {

	rows, err := e.QueryContext(ctx,
		`select coalesce(sum(n_value),0), -coalesce(sum(case when n_value < 0 then n_value end),0) from t_order_value v
		where n_user = $1 `, user)

	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		return 0, 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, 0, model.ErrorNotFound
	}
	var current, withdrawn float32
	if err := rows.Scan(&current, &withdrawn); err != nil {
		return 0, 0, err
	}
	if err = rows.Err(); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return 0, 0, err
	}
	return current, withdrawn, nil
}

func SetStatus(ctx context.Context, e model.Executor, order int, status model.AccrualStatusType) error {
	if _, err := e.ExecContext(ctx, `update t_order set s_status = $1, s_state = $2 where n_id = $3`, status, model.GetState(status), order); err != nil {
		logger.Log().Error("error add accrual", zap.Error(err))
		return err
	}
	return nil
}

func Accrual(ctx context.Context, e model.Executor, user int, order int, value float32) error {
	logger.Log().Info("ADD Accrual", zap.Int("order", order), zap.Float32("value", value), zap.Int("user", user))
	if _, err := e.ExecContext(ctx, `insert into t_order_value(n_user, n_order, n_value, dt_processed) values($1, $2, $3, current_timestamp)`, user, order, value); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return model.ErrorAlreadyExist // Игнорировать дубликат, пользователь уже существует
		}
		logger.Log().Error("error add accrual", zap.Error(err))
		return err
	}
	return nil
}

func Withdraw(ctx context.Context, e model.Executor, user int, order int, value float32) error {
	logger.Log().Info("ADD Withdraw", zap.Int("order", order), zap.Float32("value", value), zap.Int("user", user))
	if _, err := e.ExecContext(ctx, `insert into t_order_value(n_user, n_order, n_value, dt_processed) values($1, $2, $3, current_timestamp)`, user, order, -value); err != nil {
		logger.Log().Error("error insert accrual", zap.Error(err), zap.Int("order", order), zap.Float32("value", value))
		return err
	}
	return nil
}

func GetWithdrawals(ctx context.Context, e model.Executor, user int) ([]model.WithdrawnResponse, error) {

	rows, err := e.QueryContext(ctx,
		`select n_order, n_value, dt_processed from t_order_value v
		where n_value < 0 and n_user = $1
		order by dt_processed desc`, user)

	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var order int
	var sum float32
	var processed time.Time

	response := make([]model.WithdrawnResponse, 0)
	for rows.Next() {
		if err := rows.Scan(&order, &sum, &processed); err != nil {
			return nil, err
		}
		response = append(response, model.WithdrawnResponse{Order: strconv.Itoa(order), Sum: -sum, Processed: processed})
	}

	if err = rows.Err(); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}
	return response, nil
}

func GetOrderList(ctx context.Context, e model.Executor, user int) ([]model.OrderResponse, error) {

	rows, err := e.QueryContext(ctx,
		`select n_id, s_status, coalesce(sum(case when n_value>0 then n_value end),0), dt_uploaded from t_order o
		left join t_order_value v on o.n_id = v.n_order 
		where o.n_user = $1
		group by n_id, s_status, dt_uploaded
		order by dt_uploaded desc`, user)

	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var order string
	var status model.AccrualStatusType
	var accrual float32
	var uploaded time.Time

	response := make([]model.OrderResponse, 0)
	for rows.Next() {
		if err := rows.Scan(&order, &status, &accrual, &uploaded); err != nil {
			return nil, err
		}
		response = append(response, model.OrderResponse{Number: order, Status: model.GetOrderStatus(status), Accrual: accrual, Uploaded: uploaded})
	}

	if err = rows.Err(); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}
	return response, nil
}

func GetActiveOrderList(ctx context.Context, e model.Executor) ([]model.AccrualProcess, error) {

	rows, err := e.QueryContext(ctx,
		`select n_id, n_user from t_order where s_state = $1`, model.OrderStateActive)

	if err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var order int
	var user int

	response := make([]model.AccrualProcess, 0)
	for rows.Next() {
		if err := rows.Scan(&order, &user); err != nil {
			return nil, err
		}
		response = append(response, model.AccrualProcess{Order: order, User: user})
	}
	if err = rows.Err(); err != nil {
		logger.Log().Error("error", zap.Error(err))
		return nil, err
	}
	return response, nil
}
