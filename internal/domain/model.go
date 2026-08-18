package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Expense struct {
	ID          int64
	Amount      decimal.Decimal
	Category    string
	Description string
	CreatedAt   time.Time
	Budget      decimal.Decimal
}
