package main

import (
	"fmt"
	"os"
)

func main() {
	dbDsn := os.Getenv("DB_DSN")
	if dbDsn == "" {
		fmt.Println("DB_DSN environment variable is required.")
		os.Exit(1)
	}
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("PORT environment variable is required.")
		os.Exit(1)
	}

	app, err := newApplication(dbDsn, port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		fmt.Println("Expected 'migrate' or 'serve' subcommands.")
		os.Exit(1)
	}
}
