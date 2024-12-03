package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/art-es/yet-another-service/internal/core/transaction"
	"github.com/art-es/yet-another-service/internal/domain/auth"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) Add(ctx context.Context, tx transaction.Transaction, user *auth.User) error {
	sqlTx, err := getSQLTxOrBegin(tx, s.db)
	if err != nil {
		return err
	}

	const query = "INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id"

	err = sqlTx.
		QueryRowContext(ctx, query, user.Name, user.Email, user.PasswordHash).
		Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}

func (s *UserStorage) EmailExists(ctx context.Context, email string) (bool, error) {
	const query = "SELECT EXISTS(email) FROM users WHERE email=$1"

	exists := false
	err := s.db.QueryRowContext(ctx, query, email).
		Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("execute query: %w", err)
	}

	return exists, nil
}

func (s *UserStorage) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
	const query = "SELECT id, email FROM users WHERE email=$1"

	user := &auth.User{}
	err := s.db.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return user, nil
}

func (s *UserStorage) Activate(ctx context.Context, tx transaction.Transaction, userID string) error {
	sqlTx, err := getSQLTxOrBegin(tx, s.db)
	if err != nil {
		return err
	}

	const query = "UPDATE users SET activated_at=CURRENT_TIMESTAMP, updated_at=CURRENT_TIMESTAMP WHERE id=$1"

	if _, err = sqlTx.ExecContext(ctx, query, userID); err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}
