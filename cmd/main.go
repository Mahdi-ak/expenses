package main

import (
	"expenses/bootstrap"
	"expenses/handler"
	"expenses/internal/application/expense"
	"expenses/internal/domain"
	"net/http"

	postgres "expenses/internal/infrastructure/postgres"
	sqlite "expenses/internal/infrastructure/sqlite"

	"log"
)

func main() {

	config := bootstrap.LoadConfig()

	var repository domain.Repository

	switch config.DatabaseDriver {
	case "postgres":
		db, err := bootstrap.InitPostgresql()
		if err != nil {
			log.Fatal(err)
		}
		repository = postgres.NewPostgreSQLRepository(db)

	case "sqlite":
		db, err := bootstrap.InitSQLLite(config.SQLitePath)
		if err != nil {
			log.Fatal(err)
		}
		repository = sqlite.NewSQLLite(db)

	default:
		log.Fatalf("unsupported databse driver: %s", config.DatabaseDriver)

	}

	service := expense.NewService(repository)
	handler := handler.NewHandler(service)

	r := bootstrap.NewRouter(handler)

	log.Println("database:", config.DatabaseDriver)
	log.Println("server running", config.ApplicationPort)

	if err := http.ListenAndServe(config.ApplicationPort, r.Setup()); err != nil {
		log.Fatal(err)
	}

}
