package bootstrap

import (
	infrastructure "expenses/internal/infrastructure/postgres"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgresql() (*gorm.DB, error) {
	host := viper.GetString("POSTGRES_HOST")
	port := viper.GetString("POSTGRES_PORT")
	user := viper.GetString("POSTGRES_USER")
	password := viper.GetString("POSTGRES_PASSWORD")
	dbName := viper.GetString("POSTGRES_DB")
	sslMode := viper.GetString("DB_SSLMODE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host,
		user,
		password,
		dbName,
		port,
		sslMode,
	)

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&infrastructure.PostgresExpense{}); err != nil {
		return nil, err
	}

	return db, nil
}
