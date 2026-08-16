package bootstrap

import (
	"expenses/internal/expense"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB(path string) (*gorm.DB, error) {

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	err = db.AutoMigrate(&expense.Expense{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}
	return db, nil

}
