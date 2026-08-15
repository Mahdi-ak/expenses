package expense

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type SQLiteRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Create(ctx context.Context, expense Expense) (Expense, error) {

	result := r.db.WithContext(ctx).Create(&expense)
	if result.Error != nil {
		return Expense{}, result.Error
	}
	return expense, nil
}

func (r *SQLiteRepository) GetByID(ctx context.Context, id int64) (Expense, error) {

	var expense Expense

	result := r.db.WithContext(ctx).First(&expense, id)
	if result.Error != nil {
		// if result.Error == gorm.ErrRecordNotFound {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Expense{}, ErrExpenseNotFound
		}
		return Expense{}, result.Error
	}
	return expense, nil

}
func (r *SQLiteRepository) Delete(ctx context.Context, id int64) error {

	result := r.db.WithContext(ctx).Delete(&Expense{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrExpenseNotFound
	}
	return nil
}

func (r *SQLiteRepository) GetAll(ctx context.Context, filter Filter) ([]Expense, error) {

	var expenses []Expense

	query := r.db.WithContext(ctx)

	if filter.Category != nil {
		query = query.Where("category = ?", *filter.Category)
	}

	if filter.Date != nil {
		// start := filter.Date.Truncate(24 * time.Hour)
		// end := start.Add(24 * time.Hour)
		date := *filter.Date

		start := time.Date(
			date.Year(),
			date.Month(),
			date.Day(),
			0, 0, 0, 0,
			date.Location(),
		)

		end := start.AddDate(0, 0, 1)
		query = query.Where("created_at >= ? AND created_at < ?", start, end)
	}

	result := query.Find(&expenses)

	if result.Error != nil {
		return nil, result.Error
	}

	return expenses, nil
}
