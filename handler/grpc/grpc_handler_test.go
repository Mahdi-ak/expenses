package handler_test

import (
	"context"
	"testing"

	handler "expenses/handler/grpc"
	"expenses/internal/domain"
	pb "expenses/proto/pb"
	"expenses/testing/mocks"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func stringPtr(s string) *string {
	return &s
}

func int64Ptr(v int64) *int64 {
	return &v
}

func TestCreateExpense(t *testing.T) {
	service := mocks.NewMockServiceInterface(t)

	expected := domain.Expense{
		ID:     1,
		Amount: decimal.NewFromInt(100),
	}

	service.On("Create", mock.Anything, mock.Anything).
		Return(expected, nil)

	h := handler.NewGrpcHandler(service)

	req := pb.CreateExpenseRequest_builder{
		Amount:      stringPtr("100"),
		Category:    stringPtr("food"),
		Description: stringPtr("lunch"),
		Budget:      stringPtr("500"),
	}.Build()

	res, err := h.CreateExpense(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, int64(1), res.GetExpense().GetId())
	assert.Equal(t, "100", res.GetExpense().GetAmount())
}

func TestCreateExpenseInvalidAmount(t *testing.T) {
	service := mocks.NewMockServiceInterface(t)
	h := handler.NewGrpcHandler(service)

	req := pb.CreateExpenseRequest_builder{
		Amount: stringPtr("abc"),
	}.Build()

	_, err := h.CreateExpense(context.Background(), req)

	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetExpenses(t *testing.T) {
	service := mocks.NewMockServiceInterface(t)

	service.On("GetAll", mock.Anything, mock.Anything).
		Return([]domain.Expense{}, nil)

	h := handler.NewGrpcHandler(service)

	req := pb.GetExpensesRequest_builder{}.Build()

	res, err := h.GetExpenses(context.Background(), req)

	require.NoError(t, err)
	assert.Empty(t, res.GetExpenses())
}

func TestGetExpense(t *testing.T) {
	service := mocks.NewMockServiceInterface(t)

	expected := domain.Expense{
		ID:     1,
		Amount: decimal.NewFromInt(100),
	}

	service.On("GetByID", mock.Anything, int64(1)).
		Return(expected, nil)

	h := handler.NewGrpcHandler(service)

	req := pb.GetExpenseRequest_builder{
		Id: int64Ptr(1),
	}.Build()

	res, err := h.GetExpense(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, int64(1), res.GetExpense().GetId())
	assert.Equal(t, "100", res.GetExpense().GetAmount())
}

func TestGetExpenseNotFound(t *testing.T) {
	service := mocks.NewMockServiceInterface(t)

	service.On("GetByID", mock.Anything, int64(1)).
		Return(domain.Expense{}, domain.ErrExpenseNotFound)

	h := handler.NewGrpcHandler(service)

	req := pb.GetExpenseRequest_builder{
		Id: int64Ptr(1),
	}.Build()

	_, err := h.GetExpense(context.Background(), req)

	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestDeleteExpense(t *testing.T) {
	service := mocks.NewMockServiceInterface(t)

	service.On("Delete", mock.Anything, int64(1)).
		Return(nil)

	h := handler.NewGrpcHandler(service)

	req := pb.DeleteExpenseRequest_builder{
		Id: int64Ptr(1),
	}.Build()

	_, err := h.DeleteExpense(context.Background(), req)

	assert.NoError(t, err)
}
