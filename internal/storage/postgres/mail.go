package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/art-es/yet-another-service/internal/domain/shared/models"
)

type MailStorage struct {
	db *sql.DB
}

func NewMailStorage(db *sql.DB) *MailStorage {
	return &MailStorage{
		db: db,
	}
}

func (s *MailStorage) Get(ctx context.Context) ([]models.Mail, error) {
	const query = "SELECT id, address, subject, content FROM mails LIMIT 20"

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	var mails []models.Mail
	for rows.Next() {
		var mail models.Mail

		err = rows.Scan(&mail.ID, &mail.Address, &mail.Subject, &mail.Content)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		mails = append(mails, mail)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}

	return mails, nil
}

func (s *MailStorage) Save(ctx context.Context, mails []models.Mail) error {
	if len(mails) == 0 {
		return errors.New("nothing to save")
	}

	if mails[0].Stored() {
		return s.update(ctx, mails)
	}

	return s.insert(ctx, mails)
}

func (s *MailStorage) insert(ctx context.Context, mails []models.Mail) error {
	query := "INSERT INTO mails (address, subject, content) VALUES "
	args := make([]any, 0, len(mails)*3)
	for _, mail := range mails {
		index := len(args)
		query += fmt.Sprintf("($%d, $%d, $%d),", index, index+1, index+2)
		args = append(args, mail.Address, mail.Subject, mail.Content)
	}

	query = query[:len(query)-1] // remove last comma

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}

func (s *MailStorage) update(ctx context.Context, mails []models.Mail) error {
	mailedIDs := make([]string, 0, len(mails))
	for _, mail := range mails {
		if mail.Mailed {
			mailedIDs = append(mailedIDs, mail.ID)
		}
	}

	if len(mailedIDs) == 0 {
		return errors.New("nothing to update")
	}

	const query = "UPDATE mails SET mailed_at=CURRENT_TIMESTAMP WHERE id=ANY($1)"

	_, err := s.db.ExecContext(ctx, query, pq.Array(mailedIDs))
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}
