package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"
)

const limit = 21

type ArticleStorage struct {
	db *sql.DB
}

func NewArticleStorage(db *sql.DB) *ArticleStorage {
	return &ArticleStorage{db: db}
}

func (s *ArticleStorage) Get(ctx context.Context, in *dto.GetArticlesIn) (*dto.GetArticlesOut, error) {
	var args []any
	query := `SELECT slug, title, content, author_id FROM articles`

	if in.FromSlug != nil {
		query += " WHERE slug >= $1"
		args = append(args, *in.FromSlug)
	}

	query += fmt.Sprintf(" LIMIT $%d", len(args))
	args = append(args, limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query: %w", err)
	}
	defer rows.Close()

	articles := make([]dto.Article, 0, limit)
	for rows.Next() {
		var article dto.Article
		if err = rows.Scan(&article.Slug, &article.Title, &article.Content, &article.AuthorID); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	var hasMore bool
	if len(articles) >= limit {
		articles = articles[:limit]
		hasMore = true
	}

	return &dto.GetArticlesOut{
		Articles: articles,
		HasMore:  hasMore,
	}, nil
}
