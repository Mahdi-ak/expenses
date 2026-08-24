package expense

import (
	"context"
	"errors"
	"expenses/internal/domain"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidExpenseID = errors.New("expense id must be positive")
)

type Service struct {
	repository domain.Repository
}

type ExpenseSummary struct {
	TotalAmount decimal.Decimal `json:"total_amount"`
	TotalBudget decimal.Decimal `json:"total_budget"`
	Difference  decimal.Decimal `json:"difference"`
}

func NewService(repository domain.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, expense domain.Expense) (domain.Expense, error) {
	return s.repository.Create(ctx, expense)
}

func (s *Service) GetAll(ctx context.Context, filter domain.Filter) ([]domain.Expense, error) {
	return s.repository.GetAll(ctx, filter)
}

func (s *Service) Delete(ctx context.Context, id int64) error {

	if id <= 0 {
		return ErrInvalidExpenseID
	}

	return s.repository.Delete(ctx, id)
}

func (s *Service) GetByID(ctx context.Context, id int64) (domain.Expense, error) {

	if id <= 0 {
		return domain.Expense{}, ErrInvalidExpenseID
	}
	return s.repository.GetByID(ctx, id)
}

func (s *Service) GetExpensesWithSummary(ctx context.Context) ([]domain.Expense, ExpenseSummary, error) {
	exspenses, err := s.repository.GetAll(ctx, domain.Filter{})
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
