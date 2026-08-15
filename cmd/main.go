package main

import (
	"expenses/internal/database"
	"expenses/internal/expense"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	// ctx := context.Background()

	db, err := database.New("exenses.db")

	if err != nil {
		log.Fatal(err)
	}

	repository := expense.New(db)

	service := expense.NewService(repository)

	handler := expense.NewHandler(service)

	router := chi.NewRouter()

	router.Post(
		"/expenses",
		handler.Create,
	)

	router.Get(
		"/expenses",
		handler.GetAll,
	)

	router.Get(
		"/expenses/{id}",
		handler.GetByID,
	)

	router.Delete(
		"/expenses/{id}",
		handler.Delete,
	)

	log.Println("server running :8080")

	err = http.ListenAndServe(
		":8080",
		router,
	)

	if err != nil {
		log.Fatal(err)
	}

}
