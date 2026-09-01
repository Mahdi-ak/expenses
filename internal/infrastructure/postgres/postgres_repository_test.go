package infrastructure_test

import (
	"context"
	"expenses/internal/domain"
	infrastructure "expenses/internal/infrastructure/postgres"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgreSQLRepositorySuite struct {
	suite.Suite

	container *postgres.PostgresContainer
	pool      *pgxpool.Pool
	repo      domain.Repository
}

func (s *PostgreSQLRepositorySuite) SetupSuite() {
	ctx := context.Background()

	container, err := postgres.Run(
		ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("expenses_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)

	s.Require().NoError(err)

	s.container = container

	connString, err := container.ConnectionString(
		ctx,
		"sslmode=disable",
	)

	s.Require().NoError(err)

	pool, err := pgxpool.New(ctx, connString)

	s.Require().NoError(err)

	err = pool.Ping(ctx)

	s.Require().NoError(err)

	s.pool = pool

	_, err = s.pool.Exec(ctx, `
		CREATE TABLE expenses (
			id BIGSERIAL PRIMARY KEY,
			amount NUMERIC(10, 2) NOT NULL,
			category VARCHAR(100) NOT NULL,
			description TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			budget NUMERIC(10, 2) NOT NULL
		)
	`)

	s.Require().NoError(err)

	s.repo = infrastructure.NewPostgreSQLRepository(pool)
}

func (s *PostgreSQLRepositorySuite) SetupTest() {
	ctx := context.Background()

	_, err := s.pool.Exec(
		ctx,
		"TRUNCATE TABLE expenses RESTART IDENTITY CASCADE",
	)

	s.Require().NoError(err)
}

func (s *PostgreSQLRepositorySuite) TearDownSuite() {
	ctx := context.Background()

	if s.pool != nil {
		s.pool.Close()
	}

	if s.container != nil {
		err := s.container.Terminate(ctx)
		s.Require().NoError(err)
	}
}

func (s *PostgreSQLRepositorySuite) TestPostgreSQLRepository_Create() {
	ctx := context.Background()

	expense := domain.Expense{
		Amount:      decimal.NewFromFloat(150.50),
		Category:    "food",
		Description: "Lunch",
		Budget:      decimal.NewFromFloat(500),
	}

	got, err := s.repo.Create(ctx, expense)

	s.Require().NoError(err)

	s.NotZero(got.ID)

	s.True(
		expense.Amount.Equal(got.Amount),
	)

	s.Equal(
		expense.Category,
		got.Category,
	)

	s.Equal(
		expense.Description,
		got.Description,
	)

	s.True(
		expense.Budget.Equal(got.Budget),
	)

	s.False(got.CreatedAt.IsZero())
}

func (s *PostgreSQLRepositorySuite) TestPostgreSQLRepository_GetByID() {
	ctx := context.Background()

	expense := domain.Expense{
		Amount:      decimal.NewFromFloat(250.75),
		Category:    "transport",
		Description: "Taxi",
		Budget:      decimal.NewFromFloat(1000),
	}

	created, err := s.repo.Create(ctx, expense)

	s.Require().NoError(err)

	got, err := s.repo.GetByID(ctx, created.ID)

	s.Require().NoError(err)

	s.Equal(created.ID, got.ID)

	s.True(
		created.Amount.Equal(got.Amount),
	)

	s.Equal(
		created.Category,
		got.Category,
	)

	s.Equal(
		created.Description,
		got.Description,
	)

	s.True(
		created.Budget.Equal(got.Budget),
	)

	s.WithinDuration(
		created.CreatedAt,
		got.CreatedAt,
		time.Second,
	)
}

func (s *PostgreSQLRepositorySuite) TestPostgreSQLRepository_GetByID_NotFound() {
	ctx := context.Background()

	got, err := s.repo.GetByID(ctx, 999999)

	s.ErrorIs(
		err,
		domain.ErrExpenseNotFound,
	)

	s.Equal(
		domain.Expense{},
		got,
	)
}

func (s *PostgreSQLRepositorySuite) TestPostgreSQLRepository_GetAll() {
	ctx := context.Background()

	expenses := []domain.Expense{
		{
			Amount:      decimal.NewFromFloat(100),
			Category:    "food",
			Description: "Breakfast",
			Budget:      decimal.NewFromFloat(500),
		},
		{
			Amount:      decimal.NewFromFloat(200),
			Category:    "transport",
			Description: "Taxi",
			Budget:      decimal.NewFromFloat(1000),
		},
	}

	for _, expense := range expenses {
		_, err := s.repo.Create(ctx, expense)

		s.Require().NoError(err)
	}

	got, err := s.repo.GetAll(
		ctx,
		domain.Filter{},
	)

	s.Require().NoError(err)

	s.Len(got, 2)
}

func (s *PostgreSQLRepositorySuite) TestPostgreSQLRepository_GetAll_FilterByCategory() {
	ctx := context.Background()

	_, err := s.repo.Create(ctx, domain.Expense{
		Amount:      decimal.NewFromFloat(100),
		Category:    "food",
		Description: "Breakfast",
		Budget:      decimal.NewFromFloat(500),
	})

	s.Require().NoError(err)

	_, err = s.repo.Create(ctx, domain.Expense{
		Amount:      decimal.NewFromFloat(200),
		Category:    "transport",
		Description: "Taxi",
		Budget:      decimal.NewFromFloat(1000),
	})

	s.Require().NoError(err)

	category := "food"

	got, err := s.repo.GetAll(
		ctx,
		domain.Filter{
			Category: &category,
		},
	)

	s.Require().NoError(err)

	s.Len(got, 1)

	s.Equal(
		"food",
		got[0].Category,
	)

	s.True(
		decimal.NewFromFloat(100).Equal(got[0].Amount),
	)
}

func (s *PostgreSQLRepositorySuite) TestPostgreSQLRepository_GetAll_FilterByAmount() {
	ctx := context.Background()

	_, err := s.repo.Create(ctx, domain.Expense{
		Amount:      decimal.NewFromFloat(100),
		Category:    "food",
		Description: "Breakfast",
		Budget:      decimal.NewFromFloat(500),
	})

	s.Require().NoError(err)

	_, err = s.repo.Create(ctx, domain.Expense{
		Amount:      decimal.NewFromFloat(200),
		Category:    "transport",
		Description: "Taxi",
		Budget:      decimal.NewFromFloat(1000),
	})

	s.Require().NoError(err)

	amount := decimal.NewFromFloat(100)

	got, err := s.repo.GetAll(
		ctx,
		domain.Filter{
			Amount: &amount,
		},
	)

	s.Require().NoError(err)

	s.Len(got, 1)

	s.True(
		amount.Equal(got[0].Amount),
	)
}

func (s *PostgreSQLRepositorySuite) TestPostgreSQLRepository_GetAll_FilterByDate() {
	ctx := context.Background()

	filterDate := time.Now()

	_, err := s.repo.Create(ctx, domain.Expense{
		Amount:      decimal.NewFromFloat(100),
		Category:    "food",
		Description: "Breakfast",
		Budget:      decimal.NewFromFloat(500),
	})

	s.Require().NoError(err)

	_, err = s.repo.Create(ctx, domain.Expense{
		Amount:      decimal.NewFromFloat(200),
		Category:    "transport",
		Description: "Taxi",
		Budget:      decimal.NewFromFloat(1000),
	})

	s.Require().NoError(err)

	got, err := s.repo.GetAll(
		ctx,
		domain.Filter{
			Date: &filterDate,
		},
	)

	s.Require().NoError(err)

	s.Len(got, 2)

	for _, expense := range got {
		s.Equal(
			filterDate.Format("2006-01-02"),
			expense.CreatedAt.Format("2006-01-02"),
		)
	}
}

func (s *PostgreSQLRepositorySuite) TestPostgreSQLRepository_Delete() {
	ctx := context.Background()

	expense := domain.Expense{
		Amount:      decimal.NewFromFloat(300),
		Category:    "shopping",
		Description: "Shoes",
		Budget:      decimal.NewFromFloat(1000),
	}

	created, err := s.repo.Create(ctx, expense)

	s.Require().NoError(err)

	err = s.repo.Delete(
		ctx,
		created.ID,
	)

	s.Require().NoError(err)

	_, err = s.repo.GetByID(
		ctx,
		created.ID,
	)

	s.ErrorIs(
		err,
		domain.ErrExpenseNotFound,
	)
}

func (s *PostgreSQLRepositorySuite) TestPostgreSQLRepository_Delete_NotFound() {
	ctx := context.Background()

	err := s.repo.Delete(
		ctx,
		999999,
	)

	s.ErrorIs(
		err,
		domain.ErrExpenseNotFound,
	)
}

func TestPostgreSQLRepositorySuite(t *testing.T) {
	suite.Run(
		t,
		new(PostgreSQLRepositorySuite),
	)
}
