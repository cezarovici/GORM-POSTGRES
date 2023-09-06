package main

import (
	"context"
	"database/sql"
	"log"

	"time"

	"github.com/cezarovici/GORM-POSTGRES/app"
	"github.com/cezarovici/GORM-POSTGRES/infra/postgres"
)

func main() {
	db, err := sql.Open("pgx", "postgres://cezar:cezar@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := postgres.NewPostgresRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.RunRepositoryDemo(ctx, *repo)
}
