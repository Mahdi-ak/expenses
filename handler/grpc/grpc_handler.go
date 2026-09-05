package handler

import (
	"context"
	"errors"
	"time"

	"expenses/internal/application/expense"
	"expenses/internal/domain"
	pb "expenses/proto/pb"

	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	pb.UnimplementedExpenseServiceServer
	service expense.ServiceInterface
}

func NewGrpcHandler(service expense.ServiceInterface) *GrpcHandler {
	return &GrpcHandler{service: service}
}

func (s *GrpcHandler) CreateExpense(ctx context.Context, req *pb.CreateExpenseRequest) (*pb.CreateExpenseResponse, error) {
	amount, err := decimal.NewFromString(req.GetAmount())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid amount")
	}
	budget, err := decimal.NewFromString(req.GetBudget())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid budget")
	}

	exp, err := domain.NewExpense(amount, req.GetCategory(), req.GetDescription(), budget)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := s.service.Create(ctx, *exp)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return pb.CreateExpenseResponse_builder{
		Expense: toProtoExpense(result),
	}.Build(), nil
}

func (s *GrpcHandler) GetExpenses(ctx context.Context, req *pb.GetExpensesRequest) (*pb.GetExpensesResponse, error) {
	var filter domain.Filter

	if req.GetCategory() != "" {
		c := req.GetCategory()
		filter.Category = &c
	}

	if req.GetDate() != "" {
		t, err := time.Parse("2006-01-02", req.GetDate())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid date")
		}
		filter.Date = &t
	}

	if req.GetAmount() != "" {
		parsed, err := decimal.NewFromString(req.GetAmount())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid amount")
		}
		filter.Amount = &parsed
	}

	expenses, err := s.service.GetAll(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	pbExpenses := make([]*pb.Expense, 0, len(expenses))
	for _, e := range expenses {
		pbExpenses = append(pbExpenses, toProtoExpense(e))
	}

	return pb.GetExpensesResponse_builder{
		Expenses: pbExpenses,
	}.Build(), nil
}

func (s *GrpcHandler) GetExpense(ctx context.Context, req *pb.GetExpenseRequest) (*pb.GetExpenseResponse, error) {
	result, err := s.service.GetByID(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, domain.ErrExpenseNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, expense.ErrInvalidExpenseID) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return pb.GetExpenseResponse_builder{
		Expense: toProtoExpense(result),
	}.Build(), nil
}

func (s *GrpcHandler) DeleteExpense(ctx context.Context, req *pb.DeleteExpenseRequest) (*pb.DeleteExpenseResponse, error) {
	err := s.service.Delete(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, domain.ErrExpenseNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, expense.ErrInvalidExpenseID) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.DeleteExpenseResponse{}, nil
}

func (h *GrpcHandler) GetSummary(ctx context.Context, req *pb.GetSummaryRequest) (*pb.GetSummaryResponse, error) {
	expenses, summary, err := h.service.GetExpensesWithSummary(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	pbExpenses := make([]*pb.Expense, 0, len(expenses))

	for _, exp := range expenses {
		pbExpenses = append(pbExpenses, toProtoExpense(exp))
	}

	totalAmount := summary.TotalAmount.String()
	totalBudget := summary.TotalBudget.String()
	difference := summary.Difference.String()

	return pb.GetSummaryResponse_builder{
		Expenses: pbExpenses,
		Summary: pb.ExpenseSummary_builder{
			TotalAmount: &totalAmount,
			TotalBudget: &totalBudget,
			Difference:  &difference,
		}.Build(),
	}.Build(), nil
}
