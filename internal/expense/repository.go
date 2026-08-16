package expense

import (
	"context"
	"errors"
)

var ErrExpenseNotFound = errors.New("expense not found")

type Repository interface {
	Create(ctx context.Context, expense Expense) (Expense, error)

	GetAll(ctx context.Context, filter Filter) ([]Expense, error)

	GetByID(ctx context.Context, id int64) (Expense, error)

	Delete(ctx context.Context, id int64) error
}
