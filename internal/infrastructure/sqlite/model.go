package infrastructure

import (
	"time"

	"github.com/shopspring/decimal"
)

type Expense struct {
	ID          int64           `gorm:"primaryKey" json:"id"`
	Amount      decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"amount"`
	Category    string          `gorm:"not null" json:"category"`
	Description string          `gorm:"type:text" json:"description"`
	CreatedAt   time.Time       `gorm:"autoCreateTime;not null" json:"created_at"`
	Budget      decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"budget"`
}
