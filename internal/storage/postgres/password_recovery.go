package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"
	"github.com/art-es/yet-another-service/internal/core/transaction"
)

type PasswordRecoveryStorage struct {
	db *sql.DB
}

func NewPasswordRecoveryStorage(db *sql.DB) *PasswordRecoveryStorage {
	return &PasswordRecoveryStorage{
		db: db,
	}
}

func (s *PasswordRecoveryStorage) Find(ctx context.Context, token string) (*dto.PasswordRecovery, error) {
	const query = "SELECT token, user_id FROM password_recoveries WHERE token=$1"

	recovery := &dto.PasswordRecovery{}
	err := s.db.QueryRowContext(ctx, query, token).
		Scan(&recovery.Token, &recovery.UserID)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}

	return recovery, nil
}

func (s *PasswordRecoveryStorage) Delete(ctx context.Context, tx transaction.Transaction, token string) error {
	sqlTx, err := getSQLTxOrBegin(tx, s.db)
	if err != nil {
		return err
	}

	const query = "DELETE FROM password_recoveries WHERE token=$1"

	if _, err = sqlTx.ExecContext(ctx, query, token); err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}

func (s *PasswordRecoveryStorage) Save(ctx context.Context, tx transaction.Transaction, recovery *dto.PasswordRecovery) error {
	sqlTx, err := getSQLTxOrBegin(tx, s.db)
	if err != nil {
		return err
	}

	const query = "INSERT INTO password_recoveries (user_id) VALUES ($1) RETURNING token"

	err = sqlTx.QueryRowContext(ctx, query, recovery.UserID).
		Scan(&recovery.Token)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}
