package model

import (
	"errors"
	"time"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"go.uber.org/zap"
)

var ErrorAlreadyExist = errors.New("already exists")
var ErrorInvalidOrder = errors.New("invalid order")
var ErrorInvalidStatus = errors.New("invalid accrual status")
var ErrorNotFound = errors.New("not found")
var ErrorInvalidBalance = errors.New("invalid balance")
var ErrorInvalidLuhn = errors.New("invalid luhn value")
var ErrorNotOwned = errors.New("order not owned user")

/*
	{
	    "login": "<login>",
	    "password": "<password>"
	}
*/
type OrderStatusType string

const (
	OrderStatusUndefined  OrderStatusType = "UNDEFINED"
	OrderStatusNew        OrderStatusType = "NEW"        // заказ загружен в систему, но не попал в обработку;
	OrderStatusProcessing OrderStatusType = "PROCESSING" // вознаграждение за заказ рассчитывается;
	OrderStatusInvalid    OrderStatusType = "INVALID"    // система расчёта вознаграждений отказала в расчёте;
	OrderStatusProcessed  OrderStatusType = "PROCESSED"  // данные по заказу проверены и информация о расчёте успешно получена.

)

type AccrualStatusType string

const (
	AccrualStatusEmpty      AccrualStatusType = ""
	AccrualStatusRegistered AccrualStatusType = "REGISTERED"
	AccrualStatusInvalid    AccrualStatusType = "INVALID"    // система расчёта вознаграждений отказала в расчёте;
	AccrualStatusProcessing AccrualStatusType = "PROCESSING" // вознаграждение за заказ рассчитывается;
	AccrualStatusProcessed  AccrualStatusType = "PROCESSED"  // данные по заказу проверены и информация о расчёте успешно получена.
)

type OrderStateType string

const (
	OrderStateUndefined OrderStateType = "UNDEFINED"
	OrderStateActive    OrderStateType = "ACTIVE"
	OrderStateFinish    OrderStateType = "FINISH"
)

func GetOrderStatus(x AccrualStatusType) OrderStatusType {

	switch x {
	case AccrualStatusEmpty, AccrualStatusRegistered:
		return OrderStatusNew
	case AccrualStatusProcessing:
		return OrderStatusProcessing
	case AccrualStatusInvalid:
		return OrderStatusInvalid
	case AccrualStatusProcessed:
		return OrderStatusProcessed
	}
	logger.Log().Warn("OrderStatusUndefined", zap.String("accrual status", string(x)))
	return OrderStatusUndefined
}

func GetState(x AccrualStatusType) OrderStateType {
	switch x {
	case AccrualStatusRegistered, AccrualStatusProcessing:
		return OrderStateActive
	case AccrualStatusInvalid, AccrualStatusProcessed:
		return OrderStateFinish
	}
	return OrderStateUndefined
}

type AccrualProcess struct {
	Order int
	User  int
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

/*
	{
	    "number": "9278923470",
	    "status": "PROCESSED",
	    "accrual": 500,
	    "uploaded_at": "2020-12-10T15:15:45+03:00"
	}
*/
type OrderResponse struct {
	Number   string          `json:"number"`
	Status   OrderStatusType `json:"status"`
	Accrual  float32         `json:"accrual"`
	Uploaded time.Time       `json:"uploaded_at"`
}

/*
	{
	    "current": 500.5,
	    "withdrawn": 42
	}
*/
type BalanceResponse struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

/*
	{
	    "order": "2377225624",
	    "sum": 751
	}
*/
type WithdrawnRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

/*
	{
	    "order": "2377225624",
	    "sum": 500,
	    "processed_at": "2020-12-09T16:09:57+03:00"
	}
*/
type WithdrawnResponse struct {
	Order     string    `json:"order"`
	Sum       float32   `json:"sum"`
	Processed time.Time `json:"processed_at"`
}

/*
	{
	    "order": "<number>",
	    "status": "PROCESSED",
	    "accrual": 500
	}
*/
type AccrualResponse struct {
	Order   string            `json:"order"`
	Status  AccrualStatusType `json:"status"`
	Accrual float32           `json:"accrual"`
}
