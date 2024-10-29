package postgresql

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(connectionString string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	//path for docker file:///app/internal/storage/migrations
	m, err := migrate.New("file://internal/storage/migrations", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}
