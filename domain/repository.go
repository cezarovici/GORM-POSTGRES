package domain

import (
	"context"
	"errors"
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
	All(ctx context.Context, user User) (*User, error)
	GetByName(ctx context.Context, name string) (*User, error)
	Update(ctx context.Context, name string) (*User, error)
	Delete(ctx context.Context, name string) (*User, error)
}
