package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ApplicationPort  string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSslMode  string
}

func LoadConfig() Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("APPLICATION_PORT", ":8080")
	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", "5432")
	viper.SetDefault("POSTGRES_USER", "postgres")
	viper.SetDefault("POSTGRES_PASSWORD", "postgres")
	viper.SetDefault("POSTGRES_DB", "expenses")
	viper.SetDefault("POSTGRES_SSLMODE", "disable")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Print(".env file not found")
		} else {
			log.Fatalf("error reading file : %v", err)
		}

	}

	return Config{
		ApplicationPort:  viper.GetString("APPLICATION_PORT"),
		PostgresHost:     viper.GetString("POSTGRES_HOST"),
		PostgresPort:     viper.GetString("POSTGRES_PORT"),
		PostgresUser:     viper.GetString("POSTGRES_USER"),
		PostgresPassword: viper.GetString("POSTGRES_PASSWORD"),
		PostgresDB:       viper.GetString("POSTGRES_DB"),
		PostgresSslMode:  viper.GetString("POSTGRES_SSLMODE"),
	}
}
