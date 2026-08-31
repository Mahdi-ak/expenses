package infrastructure

import (
	"context"
	"errors"

	"expenses/internal/domain"
	"expenses/internal/infrastructure/postgres/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreSQLRepository struct {
	queries *sqlc.Queries
}

func NewPostgreSQLRepository(pool *pgxpool.Pool) domain.Repository {
	return &PostgreSQLRepository{
		queries: sqlc.New(pool),
	}
}

func (r *PostgreSQLRepository) Create(ctx context.Context, expense domain.Expense) (domain.Expense, error) {
	created, err := r.queries.CreateExpense(
		ctx,
		toCreateExpenseParams(expense),
	)
	if err != nil {
		return domain.Expense{}, err
	}

	result, err := toDomainExpense(created)
	if err != nil {
		return domain.Expense{}, err
	}

	return result, nil
}

func (r *PostgreSQLRepository) GetByID(ctx context.Context, id int64) (domain.Expense, error) {
	expense, err := r.queries.GetExpenseByID(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Expense{}, domain.ErrExpenseNotFound
		}

		return domain.Expense{}, err
	}

	result, err := toDomainExpense(expense)
	if err != nil {
		return domain.Expense{}, err
	}

	return result, nil
}

func (r *PostgreSQLRepository) GetAll(ctx context.Context, filter domain.Filter) ([]domain.Expense, error) {

	var (
		expenses []sqlc.Expense
		err      error
	)

	switch {
	case filter.Category != nil:
		expenses, err = r.queries.GetAllExpensesByCategory(
			ctx,
			*filter.Category,
		)

	case filter.Date != nil:
		expenses, err = r.queries.GetAllExpensesByDate(
			ctx,
			toPostgresTimestamp(*filter.Date),
		)

	case filter.Amount != nil:
		expenses, err = r.queries.GetAllExpensesByAmount(
			ctx,
			toPostgresNumeric(*filter.Amount),
		)

	default:
		expenses, err = r.queries.GetAllExpenses(ctx)
	}

	if err != nil {
		return nil, err
	}

	result := make([]domain.Expense, 0, len(expenses))

	for _, expense := range expenses {
		mapped, err := toDomainExpense(expense)
		if err != nil {
			return nil, err
		}

		result = append(result, mapped)
	}

	return result, nil
}
func (r *PostgreSQLRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.queries.GetExpenseByID(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrExpenseNotFound
		}

		return err
	}

	err = r.queries.DeleteExpense(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
