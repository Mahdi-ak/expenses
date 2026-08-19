package main

import (
	"expenses/bootstrap"
	"expenses/handler"
	"expenses/internal/application/expense"
	infrastructure "expenses/internal/infrastructure/postgres"

	// infrastructure "expenses/internal/infrastructure/sqlite"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func main() {

	// db, err := bootstrap.InitSQLLite("exenses.db")
	db, err := bootstrap.InitPostgresql()

	if err != nil {
		log.Fatal(err)
	}

	// repository := infrastructure.NewSQLLite(db)
	repository := infrastructure.NewPostgreSQLRepository(db)

	service := expense.NewService(repository)

	handler := handler.NewHandler(service)

	r := bootstrap.NewRouter(handler)

	viper.SetConfigFile(".env")

	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	port := viper.GetString("APPLICATION_PORT")

	log.Println("server running :", port)
	err = http.ListenAndServe(
		port,
		r.Setup(),
	)

	if err != nil {
		log.Fatal(err)
	}

}
