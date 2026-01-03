package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/repository/db"
	"go.uber.org/zap"
)

const MaxProcessedOrderCount = 1000
const accrualPeriod time.Duration = time.Second

type AccrualService interface {
	ProcessOrder(model.AccrualProcess)
	StartProcess(ctx context.Context)
	Close()
}

type AcuralData struct {
	baseURL         *url.URL
	conn            model.Connector
	activeOrderList chan model.AccrualProcess
}

func New(ctx context.Context, conn model.Connector, addr string) (*AcuralData, error) {
	baseURL, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	if err := conn.PingContext(ctx); err != nil {
		return nil, err
	}
	a := &AcuralData{
		baseURL:         baseURL,
		conn:            conn,
		activeOrderList: make(chan model.AccrualProcess, MaxProcessedOrderCount),
	}
	return a, nil
}

func (p *AcuralData) Close() {
	close(p.activeOrderList)
}

func (p *AcuralData) ProcessOrder(data model.AccrualProcess) {
	p.activeOrderList <- data
}

const accrualPath = "api/orders"

func parseResponceStatus(status int) (model.OrderStateType, error) {
	switch status {
	case http.StatusTooManyRequests:
		time.Sleep(accrualPeriod)
		return model.OrderStateActive, nil
	case http.StatusNoContent, http.StatusInternalServerError:
		return model.OrderStateActive, nil
	case http.StatusOK:
		return model.OrderStateFinish, nil
	default:
		return model.OrderStateFinish, model.ErrorInvalidStatus
	}

}
func parseResponseDataStatus(status model.AccrualStatusType) (model.OrderStateType, error) {
	switch status {
	case model.AccrualStatusRegistered, model.AccrualStatusProcessing:
		return model.OrderStateActive, nil
	case model.AccrualStatusInvalid:
		logger.Log().Error("invalid status from accrual")
		return model.OrderStateFinish, model.ErrorInvalidOrder
	case model.AccrualStatusProcessed:
		return model.OrderStateFinish, nil
	default:
		return model.OrderStateFinish, model.ErrorInvalidStatus
	}
}

func (p *AcuralData) request(ctx context.Context, user int, order int) (model.OrderStateType, error) {
	logger.Log().Info("accrual request start", zap.Int("order", order))
	defer logger.Log().Info("accrual request stop", zap.Int("order", order))
	joinedURL := p.baseURL.JoinPath(accrualPath, strconv.Itoa(order))

	logger.Log().Info("accrual get path", zap.String("path", joinedURL.String()))

	client := &http.Client{}
	resp, err := client.Get(joinedURL.String())

	if err != nil {
		logger.Log().Error("error read request from accrual", zap.Error(err))
		return model.OrderStateActive, err
	}
	defer resp.Body.Close()

	logger.Log().Info("accrual request status", zap.Int("order", order), zap.Int("status", resp.StatusCode))

	state, err := parseResponceStatus(resp.StatusCode)
	if err != nil {
		return state, err
	}
	if state == model.OrderStateActive {
		return state, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log().Error("error read body from accrual", zap.Error(err))
		return model.OrderStateFinish, err
	}

	var value model.AccrualResponse
	if err := json.Unmarshal(body, &value); err != nil {
		logger.Log().Error("error parse result from accrual", zap.Error(err))
		return model.OrderStateFinish, err
	}

	logger.Log().Info("accrual request data", zap.Int("order", order), zap.String("value", string(value.Status)), zap.Float32("accrual", value.Accrual))

	state, err = parseResponseDataStatus(value.Status)
	if err != nil {
		if errors.Is(err, model.ErrorInvalidOrder) {
			logger.Log().Error("invalid status from accrual", zap.Error(err))
			db.SetStatus(ctx, p.conn, order, value.Status)
		}
		return state, err
	}
	if state == model.OrderStateActive {
		db.SetStatus(ctx, p.conn, order, value.Status)
		return state, nil
	}

	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		logger.Log().Error("error start tx", zap.Error(err))
		return state, err
	}
	if err := db.SetStatus(ctx, tx, order, value.Status); err != nil {
		logger.Log().Error("error set status from accrual", zap.Error(err))
		tx.Rollback()
		return state, err
	}
	if err := db.Accrual(ctx, tx, user, order, float32(value.Accrual)); err != nil {
		logger.Log().Error("error set status from accrual", zap.Error(err))
		tx.Rollback()
		return state, err
	}
	tx.Commit()
	return state, nil

}

func (p *AcuralData) processValue(ctx context.Context, user int, order int) {

	logger.Log().Info("start accrual process value", zap.Int("order", order))
	defer logger.Log().Info("stop accrual process value", zap.Int("order", order))
	t := time.NewTicker(accrualPeriod)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			logger.Log().Info("step accrual process value", zap.Int("order", order))
			status, err := p.request(ctx, user, order)
			if err != nil {
				logger.Log().Error("request error", zap.Error(err))
			}
			if status == model.OrderStateFinish {
				return
			}
		}
	}
}

func (p *AcuralData) StartProcess(ctx context.Context) {

	logger.Log().Info("start accrual process", zap.Error(ctx.Err()))
	for {
		select {
		case <-ctx.Done():
			logger.Log().Info("stop accrual process", zap.Error(ctx.Err()))
			return
		case data := <-p.activeOrderList:
			go p.processValue(ctx, data.User, data.Order)
		}
	}

}
