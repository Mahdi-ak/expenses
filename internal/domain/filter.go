package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Filter struct {
	Category *string
	Date     *time.Time
	Amaount  *decimal.Decimal
}
