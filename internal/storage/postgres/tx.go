package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/art-es/yet-another-service/internal/core/transaction"
)

type txCtxKey struct{}

func setTxToContext(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txCtxKey{}, tx)
}

func getTxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txCtxKey{}).(*sql.Tx)
	return tx, ok
}

func getSQLTxOrBegin(tx transaction.Transaction, db *sql.DB) (*sql.Tx, error) {
	ctx := tx.Context()

	if sqlTx, ok := getTxFromContext(ctx); ok {
		return sqlTx, nil
	}

	sqlTx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	tx.WithContext(setTxToContext(ctx, sqlTx))
	return sqlTx, nil
}
