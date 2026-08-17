package domain

import (
	"errors"

	"github.com/shopspring/decimal"
)

func NewExpense(amount decimal.Decimal, category string, description string, budget decimal.Decimal) (*Expense, error) {

	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("amount most be positive")
	}
	if budget.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("budget most be positive")
	}

	return &Expense{
		Amount:      amount,
		Category:    category,
		Description: description,
		Budget:      budget,
	}, nil

}
