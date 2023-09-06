package app

import (
	"context"

	apperrors "github.com/cezarovici/GORM-POSTGRES/app/errors"
	"github.com/cezarovici/GORM-POSTGRES/infra/postgres"
	"github.com/rs/zerolog/log"
)

func RunRepositoryDemo(ctx context.Context, userRepo postgres.PostgreSqlRepo) error {
	log.Info().Msg("Inside in RunRepositoryDemo")

	if err := userRepo.Migrate(ctx); err != nil {
		return &apperrors.AppError{
			Caller:     "App",
			MethodName: "RunRepositoryDemo",
			Issue:      err,
		}
	}

	return nil
}
