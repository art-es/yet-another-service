package transaction

import (
	"context"
	"errors"
)

type Transaction interface {
	Context() context.Context
	WithContext(ctx context.Context)
	AddRollback(rollback func())
	AddCommit(commit func() error)
	Rollback()
	Commit() error
}

func New(ctx context.Context) Transaction {
	return &transaction{
		ctx: ctx,
	}
}

type transaction struct {
	ctx       context.Context
	rollbacks []func()
	commits   []func() error
}

func (t *transaction) Context() context.Context {
	return t.ctx
}

func (t *transaction) WithContext(ctx context.Context) {
	t.ctx = ctx
}

func (t *transaction) AddRollback(rollback func()) {
	t.rollbacks = append(t.rollbacks, rollback)
}

func (t *transaction) AddCommit(commit func() error) {
	t.commits = append(t.commits, commit)
}

func (t *transaction) Rollback() {
	for _, rollback := range t.rollbacks {
		rollback()
	}
}

func (t *transaction) Commit() error {
	var errs []error
	for _, commit := range t.commits {
		if err := commit(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
