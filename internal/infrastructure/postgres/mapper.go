package infrastructure

import (
	"expenses/internal/domain"
	"expenses/internal/infrastructure/postgres/sqlc"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func toPostgresNumeric(value decimal.Decimal) pgtype.Numeric {
	return pgtype.Numeric{
		Int:   value.Coefficient(),
		Exp:   value.Exponent(),
		Valid: true,
	}
}

func toPostgresTimestamp(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  value,
		Valid: true,
	}
}

func toDecimal(value pgtype.Numeric) (decimal.Decimal, error) {
	if !value.Valid {
		return decimal.Zero, fmt.Errorf("numeric value is null")
	}

	if value.NaN {
		return decimal.Zero, fmt.Errorf("numeric value is NaN")
	}

	if value.Int == nil {
		return decimal.Zero, fmt.Errorf("numeric value has no integer")
	}

	result := decimal.NewFromBigInt(value.Int, value.Exp)

	return result, nil
}

func toDomainExpense(expense sqlc.Expense) (domain.Expense, error) {
	amount, err := toDecimal(expense.Amount)
	if err != nil {
		return domain.Expense{}, fmt.Errorf("convert amount: %w", err)
	}

	budget, err := toDecimal(expense.Budget)
	if err != nil {
		return domain.Expense{}, fmt.Errorf("convert budget: %w", err)
	}

	if !expense.CreatedAt.Valid {
		return domain.Expense{}, fmt.Errorf("created_at is null")
	}

	return domain.Expense{
		ID:          expense.ID,
		Amount:      amount,
		Category:    expense.Category,
		Description: expense.Description.String,
		CreatedAt:   expense.CreatedAt.Time,
		Budget:      budget,
	}, nil
}

func toCreateExpenseParams(expense domain.Expense) sqlc.CreateExpenseParams {
	return sqlc.CreateExpenseParams{
		Amount:   toPostgresNumeric(expense.Amount),
		Category: expense.Category,
		Description: pgtype.Text{
			String: expense.Description,
			Valid:  true,
		},
		Budget: toPostgresNumeric(expense.Budget),
	}
}
