package expense

import (
	"context"
	"errors"
	"testing"

	"expenses/internal/domain"
	"expenses/internal/domain/mocks"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetExpensesWithSummary(t *testing.T) {
	ctx := context.Background()

	t.Run("returns expenses with correct summary", func(t *testing.T) {
		repo := mocks.NewMockRepository(t)
		svc := NewService(repo)

		expenses := []domain.Expense{
			{
				ID:     1,
				Amount: decimal.NewFromFloat(100.50),
				Budget: decimal.NewFromFloat(200),
			},
			{
				ID:     2,
				Amount: decimal.NewFromFloat(49.50),
				Budget: decimal.NewFromFloat(150),
			},
		}

		repo.EXPECT().GetAll(ctx, domain.Filter{}).Return(expenses, nil).Once()

		got, summary, err := svc.GetExpensesWithSummary(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expenses, got)
		assert.True(t, decimal.NewFromFloat(150).Equal(summary.TotalAmount))
		assert.True(t, decimal.NewFromFloat(350).Equal(summary.TotalBudget))
		assert.True(t, decimal.NewFromFloat(200).Equal(summary.Difference))
	})

	t.Run("returns zero summary when no expenses", func(t *testing.T) {
		repo := mocks.NewMockRepository(t)
		svc := NewService(repo)

		repo.EXPECT().GetAll(ctx, domain.Filter{}).Return(nil, nil).Once()

		got, summary, err := svc.GetExpensesWithSummary(ctx)

		assert.NoError(t, err)
		assert.Empty(t, got)
		assert.True(t, decimal.Zero.Equal(summary.TotalAmount))
		assert.True(t, decimal.Zero.Equal(summary.TotalBudget))
		assert.True(t, decimal.Zero.Equal(summary.Difference))
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repo := mocks.NewMockRepository(t)
		svc := NewService(repo)

		expectedErr := errors.New("db failure")
		repo.EXPECT().GetAll(ctx, domain.Filter{}).Return(nil, expectedErr).Once()

		got, summary, err := svc.GetExpensesWithSummary(ctx)

		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, got)
		assert.Equal(t, ExpenseSummary{}, summary)
	})
}

func TestCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("creates expense", func(t *testing.T) {
		repo := mocks.NewMockRepository(t)
		svc := NewService(repo)

		expense := domain.Expense{
			Amount:   decimal.NewFromFloat(25.99),
			Category: "food",
			Budget:   decimal.NewFromFloat(100),
		}

		created := domain.Expense{ID: 1, Amount: expense.Amount, Category: expense.Category, Budget: expense.Budget}
		repo.EXPECT().Create(ctx, mock.IsType(domain.Expense{})).Return(created, nil).Once()

		got, err := svc.Create(ctx, expense)

		assert.NoError(t, err)
		assert.Equal(t, created, got)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repo := mocks.NewMockRepository(t)
		svc := NewService(repo)

		expectedErr := errors.New("db failure")
		repo.EXPECT().Create(ctx, mock.AnythingOfType("domain.Expense")).Return(domain.Expense{}, expectedErr).Once()

		got, err := svc.Create(ctx, domain.Expense{})

		assert.ErrorIs(t, err, expectedErr)
		assert.Equal(t, domain.Expense{}, got)
	})
}

func TestGetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("returns expenses for filter", func(t *testing.T) {
		repo := mocks.NewMockRepository(t)
		svc := NewService(repo)

		filter := domain.Filter{}
		expenses := []domain.Expense{
			{ID: 1, Category: "food"},
			{ID: 2, Category: "transport"},
		}
		repo.EXPECT().GetAll(ctx, filter).Return(expenses, nil).Once()

		got, err := svc.GetAll(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, expenses, got)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repo := mocks.NewMockRepository(t)
		svc := NewService(repo)

		expectedErr := errors.New("db failure")
		repo.EXPECT().GetAll(ctx, domain.Filter{}).Return(nil, expectedErr).Once()

		got, err := svc.GetAll(ctx, domain.Filter{})

		assert.ErrorIs(t, err, expectedErr)
		assert.Nil(t, got)
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()

	t.Run("deletes expense with valid id", func(t *testing.T) {
		repo := mocks.NewMockRepository(t)
		svc := NewService(repo)

		repo.EXPECT().Delete(ctx, int64(1)).Return(nil).Once()

		err := svc.Delete(ctx, 1)

		assert.NoError(t, err)
	})

	t.Run("returns error when id is not positive", func(t *testing.T) {
		repo := mocks.NewMockRepository(t)
		svc := NewService(repo)

		for _, id := range []int64{0, -1} {
			err := svc.Delete(ctx, id)

			assert.ErrorIs(t, err, ErrInvalidExpenseID)
		}
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repo := mocks.NewMockRepository(t)
		svc := NewService(repo)

		expectedErr := errors.New("db failure")
		repo.EXPECT().Delete(ctx, int64(5)).Return(expectedErr).Once()

		err := svc.Delete(ctx, 5)

		assert.ErrorIs(t, err, expectedErr)
	})
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("returns expense with valid id", func(t *testing.T) {
		repo := mocks.NewMockRepository(t)
		svc := NewService(repo)

		expense := domain.Expense{ID: 1, Category: "food"}
		repo.EXPECT().GetByID(ctx, int64(1)).Return(expense, nil).Once()

		got, err := svc.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expense, got)
	})

	t.Run("returns error when id is not positive", func(t *testing.T) {
		repo := mocks.NewMockRepository(t)
		svc := NewService(repo)

		for _, id := range []int64{0, -3} {
			got, err := svc.GetByID(ctx, id)

			assert.ErrorIs(t, err, ErrInvalidExpenseID)
			assert.Equal(t, domain.Expense{}, got)
		}
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repo := mocks.NewMockRepository(t)
		svc := NewService(repo)

		expectedErr := errors.New("db failure")
		repo.EXPECT().GetByID(ctx, int64(9)).Return(domain.Expense{}, expectedErr).Once()

		got, err := svc.GetByID(ctx, 9)

		assert.ErrorIs(t, err, expectedErr)
		assert.Equal(t, domain.Expense{}, got)
	})
}
