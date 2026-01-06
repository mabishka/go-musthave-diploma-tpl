package handler

import (
	"context"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/service/accrual"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/service/order"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/service/user"
)

type Server struct {
	user.UserService
	order.OrderService
	accrual.AccrualService
}

func New(ctx context.Context, conn model.Connector, addr string) (*Server, error) {
	user, err := user.New(ctx, conn)
	if err != nil {
		return nil, err
	}
	order, err := order.New(ctx, conn)
	if err != nil {
		return nil, err
	}

	accrual, err := accrual.New(ctx, conn, addr)
	if err != nil {
		return nil, err
	}
	go accrual.StartProcess(ctx)

	server := &Server{
		UserService:    user,
		OrderService:   order,
		AccrualService: accrual,
	}

	list, err := server.GetActiveOrderList(ctx)
	if err != nil {
		server.Close()
		return nil, err
	}

	for _, v := range list {
		server.ProcessOrder(v)
	}

	return server, nil
}

func (s *Server) Close() {
	s.AccrualService.Close()
}
