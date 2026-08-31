package main

import (
	"context"
	"expenses/bootstrap"
	"expenses/handler"
	"expenses/internal/application/expense"
	"expenses/internal/domain"
	"net/http"

	postgres "expenses/internal/infrastructure/postgres"

	"log"
)

func main() {

	ctx := context.Background()
	config := bootstrap.LoadConfig()

	var repository domain.Repository

	if err := bootstrap.RunMigrations(); err != nil {
		log.Fatal(err)
	}

	db, err := bootstrap.InitPostgresql(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repository = postgres.NewPostgreSQLRepository(db)

	service := expense.NewService(repository)
	handler := handler.NewHandler(service)

	r := bootstrap.NewRouter(handler)

	log.Println("server running", config.ApplicationPort)

	if err := http.ListenAndServe(config.ApplicationPort, r.Setup()); err != nil {
		log.Fatal(err)
	}

}
