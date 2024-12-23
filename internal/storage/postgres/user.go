package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"
	"github.com/art-es/yet-another-service/internal/core/transaction"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
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

func (s *UserStorage) Exists(ctx context.Context, email string) (bool, error) {
	const query = "SELECT EXISTS(email) FROM users WHERE email=$1"

	exists := false
	err := s.db.QueryRowContext(ctx, query, email).
		Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("execute query: %w", err)
	}

	return exists, nil
}

func (s *UserStorage) Find(ctx context.Context, id string) (*dto.User, error) {
	const query = "SELECT id, name, email, password_hash FROM users WHERE id=$1"

	user := &dto.User{}
	err := s.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return user, nil
}

func (s *UserStorage) FindByEmail(ctx context.Context, email string) (*dto.User, error) {
	const query = "SELECT id, name, email, password_hash FROM users WHERE email=$1"

	user := &dto.User{}
	err := s.db.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return user, nil
}

func (s *UserStorage) Save(ctx context.Context, tx transaction.Transaction, user *dto.User) error {
	if !user.Stored() {
		return s.store(ctx, tx, user)
	}

	return s.update(ctx, tx, user)
}

func (s *UserStorage) update(ctx context.Context, tx transaction.Transaction, user *dto.User) error {
	sqlTx, err := getSQLTxOrBegin(tx, s.db)
	if err != nil {
		return err
	}

	query := "UPDATE users SET password_hash=$1, updated_at=CURRENT_TIMESTAMP WHERE id=$2"

	_, err = sqlTx.ExecContext(ctx, query, user.PasswordHash, user.ID)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}

func (s *UserStorage) store(ctx context.Context, tx transaction.Transaction, user *dto.User) error {
	sqlTx, err := getSQLTxOrBegin(tx, s.db)
	if err != nil {
		return err
	}

	const query = "INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id"

	_, err = sqlTx.ExecContext(ctx, query, user.PasswordHash, user.ID)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}
