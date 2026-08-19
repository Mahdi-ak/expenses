package infrastructure

import (
	"expenses/internal/application/expense"
	"expenses/internal/domain"
)

func sqliToDomain(model expense.Expense) domain.Expense {
	return domain.Expense{
		ID:          model.ID,
		Amount:      model.Amount,
		Category:    model.Category,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
		Budget:      model.Budget,
	}
}

func toSQLiModel(d domain.Expense) expense.Expense {
	return expense.Expense{
		ID:          d.ID,
		Amount:      d.Amount,
		Category:    d.Category,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
		Budget:      d.Budget,
	}
}

func sqliToDomainList(models []expense.Expense) []domain.Expense {

	expenses := make([]domain.Expense, 0, len(models))

	for _, model := range models {
		expenses = append(expenses, sqliToDomain(model))
	}

	return expenses
}
