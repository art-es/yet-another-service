package article

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/yet-another-service/internal/core/log"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"
	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
)

type articleRepository interface {
	Get(ctx context.Context, in *dto.GetArticlesIn) (*dto.GetArticlesOut, error)
}

type authorRepository interface {
	Get(ctx context.Context, authorIDs []int64) (map[int64]*dto.ArticleAuthor, error)
}

type articleCache interface {
	articleRepository
	Add(ctx context.Context, in *dto.GetArticlesIn, out *dto.GetArticlesOut) error
}

type Service struct {
	articleStorage articleRepository
	articleCache   articleCache
	authorStorage  authorRepository
	logger         log.Logger
}

func NewService(
	articleStorage articleRepository,
	articleCache articleCache,
	authorStorage authorRepository,
	logger log.Logger,
) *Service {
	return &Service{
		articleStorage: articleStorage,
		articleCache:   articleCache,
		authorStorage:  authorStorage,
		logger:         logger,
	}
}

func (s *Service) Get(ctx context.Context, in *dto.GetArticlesIn) (*dto.GetArticlesOut, error) {
	out, err := s.articleCache.Get(ctx, in)
	switch {
	case err == nil:
		return out, nil
	case errors.Is(err, apperrors.ErrNoCache):
		// need to load from storage
	default:
		return nil, fmt.Errorf("get articles from cache: %w", err)
	}

	out, err = s.articleStorage.Get(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("get articles from storage: %w", err)
	}

	if len(out.Articles) == 0 {
		// should save empty response in cache?
		return out, nil
	}

	authorMap, err := s.authorStorage.Get(ctx, getAuthorIDs(out.Articles))
	if err != nil {
		return nil, fmt.Errorf("get authors from storage: %w", err)
	}

	for _, article := range out.Articles {
		article.Author = authorMap[article.AuthorID]
	}

	if err = s.articleCache.Add(ctx, in, out); err != nil {
		s.logger.Error().Err(err).Msg("add articles to cache error")
	}

	return out, nil
}

func getAuthorIDs(articles []dto.Article) []int64 {
	out := make([]int64, 0)
	set := make(map[int64]struct{})
	for _, article := range articles {
		if _, exists := set[article.AuthorID]; !exists {
			set[article.AuthorID] = struct{}{}
			out = append(out, article.AuthorID)
		}
	}
	return out
}
