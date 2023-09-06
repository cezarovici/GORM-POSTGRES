package postgres

import (
	"context"
	"errors"

	"github.com/cezarovici/GORM-POSTGRES/domain"
)

var (
	ErrDuplicate    = errors.New("User already exists")
	ErrNoExists     = errors.New("User does not exists")
	ErrUpdateFailed = errors.New("Update failed")
	ErrDeleteFailed = errors.New("Delete failed")
)

type Repository interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context) error
	All(ctx context.Context, user domain.User) (*domain.User, error)
	GetByName(ctx context.Context, name string) (*domain.User, error)
	Update(ctx context.Context, name string) (*domain.User, error)
	Delete(ctx context.Context, name string) (*domain.User, error)
}
