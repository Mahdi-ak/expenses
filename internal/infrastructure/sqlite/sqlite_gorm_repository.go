package infrastructure

import (
	"context"
	"errors"
	"expenses/internal/application/expense"
	"expenses/internal/domain"
	"time"

	"gorm.io/gorm"
)

type SQLiteRepository struct {
	db *gorm.DB
}

func NewSQLLite(db *gorm.DB) domain.Repository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Create(ctx context.Context, expense domain.Expense) (domain.Expense, error) {

	model := toSQLiModel(expense)
	result := r.db.WithContext(ctx).Create(&model)

	if result.Error != nil {
		return sqliToDomain(model), result.Error
	}

	return sqliToDomain(model), nil
}

func (r *SQLiteRepository) GetByID(ctx context.Context, id int64) (domain.Expense, error) {

	var model expense.Expense

	result := r.db.WithContext(ctx).First(&model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.Expense{}, domain.ErrExpenseNotFound
		}
		return domain.Expense{}, result.Error
	}
	return sqliToDomain(model), nil

}

func (r *SQLiteRepository) GetAll(ctx context.Context, filter domain.Filter) ([]domain.Expense, error) {

	var model []expense.Expense

	query := r.db.WithContext(ctx)

	if filter.Category != nil {
		query = query.Where("category = ?", *filter.Category)
	}

	if filter.Date != nil {

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

	if filter.Amount != nil {
		query = query.Where("amount >= ?", *filter.Amount)
	}

	result := query.Find(&model)

	if result.Error != nil {
		return nil, result.Error
	}

	return sqliToDomainList(model), nil
}

func (r *SQLiteRepository) Delete(ctx context.Context, id int64) error {

	result := r.db.WithContext(ctx).Delete(&expense.Expense{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrExpenseNotFound
	}
	return nil
}
