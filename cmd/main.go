package main

import (
	"expenses/bootstrap"
	"expenses/handler"
	"expenses/internal/expense"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func main() {

	db, err := bootstrap.NewDB("exenses.db")

	if err != nil {
		log.Fatal(err)
	}

	repository := expense.New(db)

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
