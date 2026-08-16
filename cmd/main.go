package main

import (
	"expenses/handler"
	"expenses/internal/database"
	"expenses/internal/expense"
	"expenses/internal/router"
	"log"
	"net/http"
)

func main() {

	db, err := database.New("exenses.db")

	if err != nil {
		log.Fatal(err)
	}

	repository := expense.New(db)

	service := expense.NewService(repository)

	handler := handler.NewHandler(service)

	r := router.New(handler)
	log.Println("server running :8080")

	err = http.ListenAndServe(
		":8080",
		r.Setup(),
	)

	if err != nil {
		log.Fatal(err)
	}

}
