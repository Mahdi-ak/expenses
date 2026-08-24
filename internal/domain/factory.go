package domain

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidAmount      = errors.New("amount must be positive")
	ErrInvalidBudget      = errors.New("budget must be positive")
	ErrInvalidCategory    = errors.New("category must not be empty")
	ErrCategoryTooLong    = errors.New("category is too long")
	ErrDescriptionTooLong = errors.New("description is too long")
	ErrInvalidID          = errors.New("id must be positive")
)

const (
	MaxCategoryLength    = 100
	MaxDescriptionLength = 500
)

func NewExpense(amount decimal.Decimal, category string, description string, budget decimal.Decimal) (*Expense, error) {

	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidAmount
	}

	if budget.LessThanOrEqual(decimal.Zero) {
		return nil, ErrInvalidBudget
	}

	if category == "" {
		return nil, ErrInvalidCategory
	}

	if len(category) > MaxCategoryLength {
		return nil, ErrCategoryTooLong
	}

	if len(description) > MaxDescriptionLength {
		return nil, ErrDescriptionTooLong
	}

	return &Expense{
		Amount:      amount,
		Category:    category,
		Description: description,
		Budget:      budget,
	}, nil
}
