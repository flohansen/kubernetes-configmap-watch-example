package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/flohansen/kubernetes-configmap-watch-example/internal/app"
	"github.com/flohansen/kubernetes-configmap-watch-example/internal/product"
	"github.com/flohansen/kubernetes-configmap-watch-example/pkg/postgres"

	_ "github.com/lib/pq"
)

var (
	pgHost   = os.Getenv("PG_HOST")
	pgPort   = os.Getenv("PG_PORT")
	pgUser   = os.Getenv("PG_USER")
	pgPass   = os.Getenv("PG_PASS")
	pgDbName = os.Getenv("PG_DBNAME")
)

func main() {
	config := postgres.Config{
		Host:     pgHost,
		Port:     pgPort,
		Username: pgUser,
		Password: pgPass,
		Database: pgDbName,
	}

	db, err := sql.Open("postgres", config.Dsn())
	if err != nil {
		log.Fatal(err)
	}

	repo := product.NewPgRepo(db)
	watcher := app.NewWatcher(repo)

	if err := watcher.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
