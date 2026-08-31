package bootstrap

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

func migrationDSN() string {
	host := viper.GetString("POSTGRES_HOST")
	port := viper.GetString("POSTGRES_PORT")
	user := viper.GetString("POSTGRES_USER")
	password := viper.GetString("POSTGRES_PASSWORD")
	dbName := viper.GetString("POSTGRES_DB")
	sslMode := viper.GetString("POSTGRES_SSLMODE")

	fmt.Println("SSL MODE:", sslMode)

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user,
		password,
		host,
		port,
		dbName,
		sslMode,
	)
}

func RunMigrations() error {
	m, err := migrate.New(
		"file://internal/infrastructure/postgres/migrations",
		migrationDSN(),
	)
	if err != nil {
		return err
	}

	defer func() {
		sourceErr, databaseErr := m.Close()

		if sourceErr != nil {
			fmt.Printf("migration source close error: %v\n", sourceErr)
		}

		if databaseErr != nil {
			fmt.Printf("migration database close error: %v\n", databaseErr)
		}
	}()

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func RollbackMigration() error {
	m, err := migrate.New(
		"file://internal/infrastructure/postgres/migrations",
		migrationDSN(),
	)
	if err != nil {
		return err
	}

	defer func() {
		sourceErr, databaseErr := m.Close()

		if sourceErr != nil {
			fmt.Printf("migration source close error: %v\n", sourceErr)
		}

		if databaseErr != nil {
			fmt.Printf("migration database close error: %v\n", databaseErr)
		}
	}()

	err = m.Steps(-1)

	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
