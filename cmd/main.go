package main

import (
	"expenses/bootstrap"
	"expenses/handler"
	"expenses/internal/expense"
	"log"
	"net/http"
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
	log.Println("server running :8080")

	err = http.ListenAndServe(
		":8080",
		r.Setup(),
	)

	if err != nil {
		log.Fatal(err)
	}

}
