package domain_test

import (
	"expenses/internal/domain"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExpense(t *testing.T) {
	tests := []struct {
		name        string
		amount      decimal.Decimal
		category    string
		description string
		budget      decimal.Decimal
		expectedErr error
	}{
		{
			name:        "creates expense successfully",
			amount:      decimal.NewFromInt(100),
			category:    "food",
			description: "lunch",
			budget:      decimal.NewFromInt(150),
			expectedErr: nil,
		},
		{
			name:        "returns error when amount is zero",
			amount:      decimal.Zero,
			category:    "food",
			budget:      decimal.NewFromInt(150),
			expectedErr: domain.ErrInvalidAmount,
		},
		{
			name:        "returns error when amount is negative",
			amount:      decimal.NewFromInt(-10),
			category:    "food",
			budget:      decimal.NewFromInt(150),
			expectedErr: domain.ErrInvalidAmount,
		},
		{
			name:        "returns error when budget is zero",
			amount:      decimal.NewFromInt(100),
			category:    "food",
			budget:      decimal.Zero,
			expectedErr: domain.ErrInvalidBudget,
		},
		{
			name:        "returns error when budget is negative",
			amount:      decimal.NewFromInt(100),
			category:    "food",
			budget:      decimal.NewFromInt(-10),
			expectedErr: domain.ErrInvalidBudget,
		},
		{
			name:        "returns error when category is empty",
			amount:      decimal.NewFromInt(100),
			category:    "",
			budget:      decimal.NewFromInt(150),
			expectedErr: domain.ErrInvalidCategory,
		},
		{
			name:        "returns error when category is too long",
			amount:      decimal.NewFromInt(100),
			category:    strings.Repeat("a", domain.MaxCategoryLength+1),
			budget:      decimal.NewFromInt(150),
			expectedErr: domain.ErrCategoryTooLong,
		},
		{
			name:        "returns error when description is too long",
			amount:      decimal.NewFromInt(100),
			category:    "food",
			description: strings.Repeat("a", domain.MaxDescriptionLength+1),
			budget:      decimal.NewFromInt(150),
			expectedErr: domain.ErrDescriptionTooLong,
		},
		{
			name:        "accepts category at maximum length",
			amount:      decimal.NewFromInt(100),
			category:    strings.Repeat("a", domain.MaxCategoryLength),
			budget:      decimal.NewFromInt(150),
			expectedErr: nil,
		},
		{
			name:        "accepts description at maximum length",
			amount:      decimal.NewFromInt(100),
			category:    "food",
			description: strings.Repeat("a", domain.MaxDescriptionLength),
			budget:      decimal.NewFromInt(150),
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			expense, err := domain.NewExpense(
				tc.amount,
				tc.category,
				tc.description,
				tc.budget,
			)

			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, expense)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, expense)

			assert.True(t, tc.amount.Equal(expense.Amount))
			assert.Equal(t, tc.category, expense.Category)
			assert.Equal(t, tc.description, expense.Description)
			assert.True(t, tc.budget.Equal(expense.Budget))
		})
	}
}
