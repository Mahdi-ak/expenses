package expense

import "context"

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, expense Expense) (Expense, error) {
	return s.repository.Create(ctx, expense)
}

func (s *Service) GetAll(ctx context.Context, filter Filter) ([]Expense, error) {
	return s.repository.GetAll(ctx, filter)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}

func (s *Service) GetByID(ctx context.Context, id int64) (Expense, error) {
	return s.repository.GetByID(ctx, id)
}
