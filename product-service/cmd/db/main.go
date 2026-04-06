package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/azar-writes-code/oolio-products-backend/pkg/app"
)

func main() {
	if len(os.Args) < 2 {
		slog.Error("subcommand is missing, use 'help' for more information")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "migrate":
		app.Migrate()
	case "drop":
		app.Drop()
	case "help":
		fmt.Println("Usage: go run cmd/db/main.go <subcommand>")
		fmt.Println("Subcommands:")
		fmt.Println("  migrate: apply database migrations")
		fmt.Println("  drop: drop database tables")
		fmt.Println("  help: show this help message")
		os.Exit(0)
	default:
		slog.Error("unknown subcommand", "subcommand", os.Args[1])
		os.Exit(1)
	}
}
