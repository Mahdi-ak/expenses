package domain_test

import (
	"expenses/internal/domain"
	"testing"

	"github.com/shopspring/decimal"
)

func TestNewExpense(t *testing.T) {
	tests := []struct {
		amount      decimal.Decimal
		category    string
		budget      decimal.Decimal
		expectError bool
	}{
		{
			amount:      decimal.NewFromInt(100),
			category:    "",
			budget:      decimal.NewFromInt(150),
			expectError: false,
		},
		{
			amount:      decimal.NewFromInt(100),
			category:    "",
			budget:      decimal.NewFromInt(-10),
			expectError: true,
		},
	}

	for _, tc := range tests {

		_, err := domain.NewExpense(
			tc.amount,
			tc.category,
			"",
			tc.budget,
		)
		if tc.expectError && err == nil {
			t.Error("expected error")
		}
		if !tc.expectError && err != nil {
			t.Error("expected error")
		}
	}

}
