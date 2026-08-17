package expense

import (
	"context"
	"errors"
	"expenses/internal/domain"
)

var ErrExpenseNotFound = errors.New("expense not found")

type Repository interface {
	Create(ctx context.Context, expense domain.Expense) (domain.Expense, error)

	GetAll(ctx context.Context, filter Filter) ([]domain.Expense, error)

	GetByID(ctx context.Context, id int64) (domain.Expense, error)

	Delete(ctx context.Context, id int64) error
}
