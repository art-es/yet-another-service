package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/art-es/yet-another-service/internal/domain/auth"

	"github.com/art-es/yet-another-service/internal/core/transaction"
)

type UserActivationStorage struct {
	db *sql.DB
}

func NewUserActivationStorage(db *sql.DB) *UserActivationStorage {
	return &UserActivationStorage{
		db: db,
	}
}

func (s *UserActivationStorage) Create(ctx context.Context, tx transaction.Transaction, userID string) (*auth.Activation, error) {
	sqlTx, err := getSQLTxOrBegin(tx, s.db)
	if err != nil {
		return nil, err
	}

	const query = "INSERT INTO user_activations (user_id) VALUES ($1) RETURNING (token, user_id)"

	activation := &auth.Activation{}
	err = sqlTx.QueryRowContext(ctx, query, userID).
		Scan(&activation.Token, &activation.UserID)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return activation, nil
}

func (s *UserActivationStorage) Delete(ctx context.Context, tx transaction.Transaction, token string) error {
	sqlTx, err := getSQLTxOrBegin(tx, s.db)
	if err != nil {
		return err
	}

	const query = "DELETE FROM user_activations WHERE token=$1"

	if _, err = sqlTx.ExecContext(ctx, query, token); err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}

func (s *UserActivationStorage) FindByToken(ctx context.Context, token string) (*auth.Activation, error) {
	const query = "SELECT token, user_id FROM user_activations WHERE token=$1"

	activation := &auth.Activation{}
	err := s.db.QueryRowContext(ctx, query, token).
		Scan(&activation.Token, &activation.UserID)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return activation, nil
}
