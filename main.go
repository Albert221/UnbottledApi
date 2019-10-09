package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	dbDsn := os.Getenv("DB_DSN")
	if dbDsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	app, err := newApplication(dbDsn)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		fmt.Println("Expected 'migrate' or 'serve' subcommands.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "migrate":
		app.Migrate()
	case "serve":
		if err := app.Serve(); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Expected 'migrate' or 'serve' subcommands.")
		os.Exit(1)
	}
}
