package expense

import (
	"context"
	"expenses/internal/domain"

	"github.com/shopspring/decimal"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, expense domain.Expense) (domain.Expense, error) {
	return s.repository.Create(ctx, expense)
}

func (s *Service) GetAll(ctx context.Context, filter Filter) ([]domain.Expense, error) {
	return s.repository.GetAll(ctx, filter)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}

func (s *Service) GetByID(ctx context.Context, id int64) (domain.Expense, error) {
	return s.repository.GetByID(ctx, id)
}

type ExpenseSummary struct {
	TotalAmount decimal.Decimal `json:"total_amount"`
	TotalBudget decimal.Decimal `json:"total_budget"`
	Difference  decimal.Decimal `json:"difference"`
}

func (s *Service) GetExpensesWithSummary(ctx context.Context) ([]domain.Expense, ExpenseSummary, error) {
	exspenses, err := s.repository.GetAll(ctx, Filter{})
	if err != nil {
		return nil, ExpenseSummary{}, err
	}

	totalAmount := decimal.Zero
	totalBudget := decimal.Zero

	for _, exspense := range exspenses {
		totalAmount = totalAmount.Add(exspense.Amount)
		totalBudget = totalBudget.Add(exspense.Budget)
	}

	summary := ExpenseSummary{
		TotalAmount: totalAmount,
		TotalBudget: totalBudget,
		Difference:  totalBudget.Sub(totalAmount),
	}
	return exspenses, summary, nil

}
