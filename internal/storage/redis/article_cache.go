package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/art-es/yet-another-service/internal/core/log"

	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"

	"github.com/redis/go-redis/v9"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"
)

type articleCacheElement struct {
	in  *dto.GetArticlesIn
	out *dto.GetArticlesOut
}

type ArticleCache struct {
	db     *redis.Client
	logger log.Logger

	cacheTimeout  time.Duration
	enrichTimeout time.Duration
	elements      chan articleCacheElement
}

func NewArticleCache(
	db *redis.Client,
	logger log.Logger,
	cacheTimeout time.Duration,
	enrichTimeout time.Duration,
) *ArticleCache {
	return &ArticleCache{
		db:            db,
		logger:        logger,
		cacheTimeout:  cacheTimeout,
		enrichTimeout: enrichTimeout,
		elements:      make(chan articleCacheElement, 5),
	}
}

func (c *ArticleCache) Get(ctx context.Context, in *dto.GetArticlesIn) (*dto.GetArticlesOut, error) {
	b, err := c.db.Get(ctx, c.key(in)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, apperrors.ErrNoCache
		}

		return nil, fmt.Errorf("execute command: %w", err)
	}

	out := &dto.GetArticlesOut{}
	if err = json.Unmarshal(b, out); err != nil {
		return nil, fmt.Errorf("unmarshal data: %w", err)
	}

	return out, nil
}

func (c *ArticleCache) Add(ctx context.Context, in *dto.GetArticlesIn, out *dto.GetArticlesOut) error {
	element := articleCacheElement{in, out}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.elements <- element:
		return nil
	}
}

func (c *ArticleCache) RunEnricher(ctx context.Context) {
	for element := range c.elements {
		if err := c.enrich(ctx, element.in, element.out); err != nil {
			c.logger.Error().Err(err).Msg("enrich article cache in redis")
		}
	}
}

func (c *ArticleCache) enrich(ctx context.Context, in *dto.GetArticlesIn, out *dto.GetArticlesOut) error {
	enrichCtx, cancel := context.WithTimeout(ctx, c.enrichTimeout)
	defer cancel()

	data, err := json.Marshal(out)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	err = c.db.Set(enrichCtx, c.key(in), data, c.cacheTimeout).Err()
	if err != nil {
		return fmt.Errorf("set data: %w", err)
	}

	return nil
}

func (c *ArticleCache) key(in *dto.GetArticlesIn) string {
	vals := url.Values{}
	if in.FromSlug != nil {
		vals.Add("from_slug", *in.FromSlug)
	}

	return "article_query:" + vals.Encode()
}
