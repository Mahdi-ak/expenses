package infrastructure

import (
	"context"
	"errors"
	"expenses/internal/domain"
	"time"

	"gorm.io/gorm"
)

type PostgreSQLRepository struct {
	db *gorm.DB
}

func NewPostgreSQLRepository(db *gorm.DB) domain.Repository {
	return &PostgreSQLRepository{
		db: db,
	}
}

func (r *PostgreSQLRepository) Create(ctx context.Context, expense domain.Expense) (domain.Expense, error) {

	model := toPostgresExpense(expense)

	result := r.db.
		WithContext(ctx).
		Create(&model)

	if result.Error != nil {
		return domain.Expense{}, result.Error
	}

	return postgresToDomainExpense(model), nil
}

func (r *PostgreSQLRepository) GetAll(ctx context.Context, filter domain.Filter) ([]domain.Expense, error) {

	var models []PostgresExpense

	query := r.db.WithContext(ctx)

	if filter.Category != nil {
		query = query.Where("category = ?", *filter.Category)
	}

	if filter.Date != nil {
		start := filter.Date.Truncate(24 * time.Hour)
		end := start.Add(24 * time.Hour)

		query = query.Where(
			"created_at >= ? AND created_at < ?",
			start,
			end,
		)
	}

	result := query.Find(&models)

	if result.Error != nil {
		return nil, result.Error
	}

	expenses := make([]domain.Expense, 0, len(models))

	for _, model := range models {
		expenses = append(expenses, postgresToDomainExpense(model))
	}

	return expenses, nil
}

func (r *PostgreSQLRepository) GetByID(ctx context.Context, id int64) (domain.Expense, error) {

	var model PostgresExpense

	result := r.db.
		WithContext(ctx).
		First(&model, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.Expense{}, domain.ErrExpenseNotFound
	}

	if result.Error != nil {
		return domain.Expense{}, result.Error
	}

	return postgresToDomainExpense(model), nil
}

func (r *PostgreSQLRepository) Delete(ctx context.Context, id int64) error {

	result := r.db.
		WithContext(ctx).
		Delete(&PostgresExpense{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrExpenseNotFound
	}

	return nil
}
