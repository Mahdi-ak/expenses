package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ApplicationPort  string
	DatabaseDriver   string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	SQLitePath       string
}

func LoadConfig() Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("APPLICATION_PORT", ":8080")
	viper.SetDefault("DATABASE_DRIVER", "sqlite")
	viper.SetDefault("SQLITE_PATH", "expenses.db")
	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", "5432")
	viper.SetDefault("POSTGRES_USER", "postgres")
	viper.SetDefault("POSTGRES_PASSWORD", "postgres")
	viper.SetDefault("POSTGRES_DB", "expenses")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Print(".env file not found")
		} else {
			log.Fatalf("error reading file : %v", err)
		}

	}

	return Config{
		ApplicationPort:  viper.GetString("APPLICATION_PORT"),
		DatabaseDriver:   viper.GetString("DATABASE_DRIVER"),
		PostgresHost:     viper.GetString("POSTGRES_HOST"),
		PostgresPort:     viper.GetString("POSTGRES_PORT"),
		PostgresUser:     viper.GetString("POSTGRES_USER"),
		PostgresPassword: viper.GetString("POSTGRES_PASSWORD"),
		PostgresDB:       viper.GetString("POSTGRES_DB"),
		SQLitePath:       viper.GetString("SQLITE_PATH"),
	}
}
