package expense

import "expenses/internal/domain"

func toDomain(model Expense) domain.Expense {
	return domain.Expense{
		ID:          model.ID,
		Amount:      model.Amount,
		Category:    model.Category,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
		Budget:      model.Budget,
	}
}

func toDBModel(expense domain.Expense) Expense {
	return Expense{
		ID:          expense.ID,
		Amount:      expense.Amount,
		Category:    expense.Category,
		Description: expense.Description,
		CreatedAt:   expense.CreatedAt,
		Budget:      expense.Budget,
	}
}

func toDomainList(models []Expense) []domain.Expense {

	expenses := make([]domain.Expense, 0, len(models))

	for _, model := range models {
		expenses = append(expenses, toDomain(model))
	}

	return expenses
}
