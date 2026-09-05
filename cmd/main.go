package main

import (
	"context"
	"expenses/bootstrap"
	"expenses/internal/application/expense"
	"expenses/internal/domain"
	"flag"

	postgres "expenses/internal/infrastructure/postgres"

	"log"
)

func main() {
	transport := flag.String("transport", "both", "transport to run: http, grpc, both")
	httpPortFlag := flag.String("http-port", "", "override APPLICATION_PORT (e.g. :8080)")
	grpcPortFlag := flag.String("grpc-port", "", "override GRPC_PORT (e.g. :50051)")
	flag.Parse()

	ctx := context.Background()
	config := bootstrap.LoadConfig()

	if *httpPortFlag != "" {
		config.ApplicationPort = *httpPortFlag
	}
	if *grpcPortFlag != "" {
		config.GrpcPort = *grpcPortFlag
	}

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

	switch *transport {
	case "http":
		bootstrap.StartHTTPServer(config.ApplicationPort, service)
	case "grpc":
		bootstrap.StartGRPCServer(config.GrpcPort, service)
	case "both":
		go bootstrap.StartHTTPServer(config.ApplicationPort, service)
		bootstrap.StartGRPCServer(config.GrpcPort, service)
	default:
		log.Fatalf("invalid -transport value %q: must be http, grpc, or both", *transport)
	}
}
