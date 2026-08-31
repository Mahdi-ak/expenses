package main

import (
	"log"
	"os"

	"expenses/bootstrap"
)

func main() {
	bootstrap.LoadConfig()
	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./cmd/migrate [up|down]")
	}

	switch os.Args[1] {
	case "up":
		if err := bootstrap.RunMigrations(); err != nil {
			log.Fatal(err)
		}

		log.Println("migrations applied successfully")

	case "down":
		if err := bootstrap.RollbackMigration(); err != nil {
			log.Fatal(err)
		}

		log.Println("migration rollback successfully")

	default:
		log.Fatalf("unknown command: %s. use up or down", os.Args[1])
	}
}
