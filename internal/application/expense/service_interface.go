package expense

import (
	"context"
	"expenses/internal/domain"
)

type ServiceInterface interface {
	Create(ctx context.Context, expense domain.Expense) (domain.Expense, error)
	GetAll(ctx context.Context, filter domain.Filter) ([]domain.Expense, error)
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (domain.Expense, error)
	GetExpensesWithSummary(ctx context.Context) ([]domain.Expense, ExpenseSummary, error)
}
