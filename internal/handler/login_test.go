package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mabishka/go-musthave-diploma-tpl/internal/logger"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/model"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/service/accrual"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/service/order"
	"github.com/mabishka/go-musthave-diploma-tpl/internal/service/user"
	"github.com/mabishka/go-musthave-diploma-tpl/mock"
	"github.com/stretchr/testify/assert"
)

func TestServer_HandlerPostLogin(t *testing.T) {

	if err := logger.InitLogger("info"); err != nil {
		panic(err)
	}

	ctx := context.Background()
	conn := &mock.Connector{}
	user, err := user.New(ctx, conn)
	assert.NoError(t, err)
	order, err := order.New(ctx, conn)
	assert.NoError(t, err)
	accrual, err := accrual.New(ctx, conn, "")
	assert.NoError(t, err)

	s := &Server{
		UserService:    user,
		OrderService:   order,
		AccrualService: accrual,
	}

	tests := []struct {
		name       string // description of this test case
		addr       string
		body       string
		method     string
		content    string
		wantStatus int
	}{
		{
			name:       "positive",
			body:       `{ "login": "user", "password": "password"}`,
			method:     http.MethodPost,
			content:    model.ContentTypeJSON,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "negative content",
			body:       `{ "login": "user", "password": "password"}`,
			method:     http.MethodPost,
			content:    model.ContentTypeHTML,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "negative method",
			body:       `{ "login": "user", "password": "password"}`,
			method:     http.MethodGet,
			content:    model.ContentTypeJSON,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "negative body",
			body:       `{ "login": "user", "password": "password}`,
			method:     http.MethodPost,
			content:    model.ContentTypeJSON,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, "/api/user/login", strings.NewReader(test.body))
			r.Header.Set(model.HeaderContentType, test.content)

			w := httptest.NewRecorder()
			s.HandlerPostLogin(w, r)

			assert.Equal(t, w.Code, test.wantStatus)
		})
	}
}
