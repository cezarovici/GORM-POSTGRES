package app

import (
	"context"
	"fmt"
	"log"

	"github.com/cezarovici/GORM-POSTGRES/infra/postgres"
)

func RunRepositoryDemo(ctx context.Context, userRepo postgres.PostgreSqlRepo) {
	fmt.Println("1. MIGRATE REPOSITORY")

	if err := userRepo.Migrate(ctx); err != nil {
		log.Fatal(err)
	}
}
