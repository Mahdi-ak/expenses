package expense

import "time"

type Expense struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	Amount      float64   `gorm:"not null" json:"amount"`
	Category    string    `gorm:"not null" json:"category"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime;not null" json:"created_at"`
}
