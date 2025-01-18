package postgres

import (
	"context"
	"database/sql"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"
)

type ArticleAuthorStorage struct {
	db *sql.DB
}

func NewArticleAuthorStorage(db *sql.DB) *ArticleAuthorStorage {
	return &ArticleAuthorStorage{db: db}
}

func (s *ArticleAuthorStorage) Get(ctx context.Context, authorIDs []int64) (map[int64]*dto.ArticleAuthor, error) {
	return nil, nil
}
