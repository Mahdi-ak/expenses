package infrastructure

import "expenses/internal/domain"

func toPostgresExpense(expense domain.Expense) PostgresExpense {
	return PostgresExpense{
		ID:          expense.ID,
		Amount:      expense.Amount,
		Category:    expense.Category,
		Description: expense.Description,
		CreatedAt:   expense.CreatedAt,
		Budget:      expense.Budget,
	}
}

func postgresToDomainExpense(expense PostgresExpense) domain.Expense {
	return domain.Expense{
		ID:          expense.ID,
		Amount:      expense.Amount,
		Category:    expense.Category,
		Description: expense.Description,
		CreatedAt:   expense.CreatedAt,
		Budget:      expense.Budget,
	}
}
