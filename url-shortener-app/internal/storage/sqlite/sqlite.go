package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/alishhh/url-shortener-app/internal/config"
	"github.com/alishhh/url-shortener-app/internal/storage"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Alias struct {
	Alias string
	URL   string
}

type Storage struct {
	logger *slog.Logger
	db     *sql.DB
}

func New(cfg *config.Config, logger *slog.Logger) *Storage {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		log.Fatalf("%s", err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		log.Fatalf("%s", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatalf("%s", err)
	}

	return &Storage{
		logger: logger,
		db:     db,
	}
}

func (s *Storage) SaveURL(ctx context.Context, urlToSave string, alias string) error {
	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO url(url, alias) values(?, ?)")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	_, err = stmt.ExecContext(ctx, urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%w", storage.ErrURLExists)
		}

		return fmt.Errorf("%w", err)
	}

	return nil
}

func (s *Storage) RemoveURL(ctx context.Context, alias string) error {
	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	_, err = stmt.ExecContext(ctx, alias)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (s *Storage) GetURLs(ctx context.Context) ([]Alias, error) {
	stmt, err := s.db.PrepareContext(ctx, "SELECT alias, url FROM url")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrURLNotFound
		}

		return nil, fmt.Errorf("%w", err)
	}

	var as []Alias
	for rows.Next() {
		var a Alias
		if err := rows.Scan(&a.Alias, &a.URL); err != nil {
			return as, err
		}
		as = append(as, a)
	}

	return as, nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	stmt, err := s.db.PrepareContext(ctx, "SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	var res string
	err = stmt.QueryRowContext(ctx, alias).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%w", err)
	}

	return res, nil
}
