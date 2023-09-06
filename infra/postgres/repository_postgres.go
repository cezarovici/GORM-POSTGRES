package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cezarovici/GORM-POSTGRES/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgreSqlRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgreSqlRepo {
	return &PostgreSqlRepo{
		db: db,
	}
}

func (r *PostgreSqlRepo) Migrate(ctx context.Context) error {
	query := `
    CREATE TABLE IF NOT EXISTS domain.Users(
        id SERIAL PRIMARY KEY,
		rank INT NOT NULL
        first_name TEXT NOT NULL,
        last_name TEXT NOT NULL,
    );`

	_, errQueryExec := r.db.ExecContext(ctx, query)

	return errQueryExec
}

func (r *PostgreSqlRepo) Create(ctx context.Context, user domain.User) (*domain.User, error) {
	var id int32

	err := r.db.QueryRowContext(ctx, "INSERT INTO users(rank,first_name, last_name) values($1, $2, $3) RETURNING id", user.Rank, user.FirstName, user.LastName).Scan(&id)

	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	user.ID = id

	return &user, nil
}

func (r *PostgreSqlRepo) All(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM domain.Users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Rank, &user.FirstName, &user.LastName); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *PostgreSqlRepo) GetByName(ctx context.Context, name string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM userss WHERE name = $1", name)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Rank, &user.FirstName, &user.LastName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoExists
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgreSqlRepo) Update(ctx context.Context, id int64, updated domain.User) (*domain.User, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE domain.Users SET rank = $1, first_name = $2 last_name = $3 WHERE id = $4", updated.Rank, updated.FirstName, updated.LastName, id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updated, nil
}

func (r *PostgreSqlRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM domain.Users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
