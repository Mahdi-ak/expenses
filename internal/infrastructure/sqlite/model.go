package infrastructure

import (
	"time"

	"github.com/shopspring/decimal"
)

type SQLiExpense struct {
	ID          int64           `gorm:"primaryKey"`
	Amount      decimal.Decimal `gorm:"type:numeric(10,2);not null"`
	Category    string          `gorm:"not null"`
	Description string          `gorm:"type:text"`
	CreatedAt   time.Time       `gorm:"autoCreateTime;not null"`
	Budget      decimal.Decimal `gorm:"type:numeric(10,2);not null"`
}
