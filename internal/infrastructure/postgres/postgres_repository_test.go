package infrastructure_test

import (
	"context"
	"expenses/internal/domain"
	infrastructure "expenses/internal/infrastructure/postgres"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	pggorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQLRepositorySuite struct {
	suite.Suite
	container *postgres.PostgresContainer
	db        *gorm.DB
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
		postgres.BasicWaitStrategies())

	s.Require().NoError(err)

	s.container = container

	conn, err := container.ConnectionString(
		ctx,
		"sslmode=disable",
	)

	s.Require().NoError(err)

	db, err := gorm.Open(pggorm.Open(conn), &gorm.Config{})

	s.Require().NoError(err)

	err = db.AutoMigrate(&infrastructure.PostgresExpense{})

	s.Require().NoError(err)

	s.db = db
	s.repo = infrastructure.NewPostgreSQLRepository(db)
}

func (s *PostgreSQLRepositorySuite) SetupTest() {
	err := s.db.Exec(
		"TRUNCATE TABLE postgres_expenses RESTART IDENTITY CASCADE",
	).Error

	s.Require().NoError(err)
}

func (s *PostgreSQLRepositorySuite) TearDownSuite() {
	ctx := context.Background()

	err := s.container.Terminate(ctx)

	s.Require().NoError(err)
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
	s.True(expense.Amount.Equal(got.Amount))
	s.Equal(expense.Category, got.Category)
	s.Equal(expense.Description, got.Description)
	s.True(expense.Budget.Equal(got.Budget))
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
	s.True(created.Amount.Equal(got.Amount))
	s.Equal(created.Category, got.Category)
	s.Equal(created.Description, got.Description)
	s.True(created.Budget.Equal(got.Budget))

	s.WithinDuration(
		created.CreatedAt,
		got.CreatedAt,
		time.Second,
	)
}

func (s *PostgreSQLRepositorySuite) TestPostgreSQLRepository_GetByID_NotFound() {
	ctx := context.Background()

	got, err := s.repo.GetByID(ctx, 999999)

	s.ErrorIs(err, domain.ErrExpenseNotFound)
	s.Equal(domain.Expense{}, got)
}

func (s *PostgreSQLRepositorySuite) TestPostgreSQLRepository_TestGetAll() {
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

	got, err := s.repo.GetAll(ctx, domain.Filter{})

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
	s.Equal("food", got[0].Category)
	s.True(
		decimal.NewFromFloat(100).Equal(got[0].Amount),
	)
}
func (s *PostgreSQLRepositorySuite) TestPostgreSQLRepository_TestDelete() {
	ctx := context.Background()

	expense := domain.Expense{
		Amount:      decimal.NewFromFloat(300),
		Category:    "shopping",
		Description: "Shoes",
		Budget:      decimal.NewFromFloat(1000),
	}

	created, err := s.repo.Create(ctx, expense)

	s.Require().NoError(err)

	err = s.repo.Delete(ctx, created.ID)

	s.Require().NoError(err)

	_, err = s.repo.GetByID(ctx, created.ID)

	s.ErrorIs(err, domain.ErrExpenseNotFound)
}
func TestPostgreSQLRepositorySuite(t *testing.T) {
	suite.Run(t, new(PostgreSQLRepositorySuite))
}
