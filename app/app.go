package app

import (
	"context"
	"fmt"
	"log"

	"github.com/cezarovici/GORM-POSTGRES/domain"
)

func RunRepositoryDemo(ctx context.Context, userRepo domain.Repository) {
	fmt.Println("1. MIGRATE REPOSITORY")

	if err := userRepo.Migrate(ctx); err != nil {
		log.Fatal(err)
	}
}
