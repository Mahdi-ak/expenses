package handler_test

import (
	"bytes"
	"context"

	handler "expenses/handler/http"
	"expenses/internal/application/expense"
	"expenses/testing/mocks"

	"expenses/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	service := mocks.NewMockServiceInterface(t)

	expected := domain.Expense{
		ID:          1,
		Amount:      decimal.NewFromInt(100),
		Category:    "food",
		Description: "lunch",
		Budget:      decimal.NewFromInt(500),
	}

	service.
		On("Create", mock.Anything, mock.Anything).
		Return(expected, nil)

	h := handler.NewHandler(service)

	body := `{
		"amount": "100",
		"category": "food",
		"description": "lunch",
		"budget": "500"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/expenses",
		bytes.NewBufferString(body),
	)

	rec := httptest.NewRecorder()

	h.Create(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestGetAll(t *testing.T) {
	service := mocks.NewMockServiceInterface(t)

	service.
		On("GetAll", mock.Anything, domain.Filter{}).
		Return([]domain.Expense{}, nil)

	h := handler.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/expenses",
		nil,
	)

	rec := httptest.NewRecorder()

	h.GetAll(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetByID(t *testing.T) {
	service := mocks.NewMockServiceInterface(t)

	expected := domain.Expense{
		ID:       1,
		Amount:   decimal.NewFromInt(100),
		Category: "food",
	}

	service.
		On("GetByID", mock.Anything, int64(1)).
		Return(expected, nil)

	h := handler.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/expenses/1",
		nil,
	)

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			routeCtx,
		),
	)

	rec := httptest.NewRecorder()

	h.GetByID(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDelete(t *testing.T) {
	service := mocks.NewMockServiceInterface(t)

	service.
		On("Delete", mock.Anything, int64(1)).
		Return(nil)

	h := handler.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/expenses/1",
		nil,
	)

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(
			req.Context(),
			chi.RouteCtxKey,
			routeCtx,
		),
	)

	rec := httptest.NewRecorder()

	h.Delete(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestGetSummary(t *testing.T) {
	service := mocks.NewMockServiceInterface(t)

	expenses := []domain.Expense{
		{
			ID:       1,
			Amount:   decimal.NewFromInt(100),
			Category: "food",
		},
	}

	summary := expense.ExpenseSummary{
		TotalAmount: decimal.NewFromInt(100),
		TotalBudget: decimal.NewFromInt(500),
		Difference:  decimal.NewFromInt(400),
	}

	service.
		On("GetExpensesWithSummary", mock.Anything).
		Return(expenses, summary, nil)

	h := handler.NewHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/expenses/summary",
		nil,
	)

	rec := httptest.NewRecorder()

	h.GetSummary(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
