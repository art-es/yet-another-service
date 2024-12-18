package redis

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthTokenBlackListStorage struct {
	db *redis.Client
}

func NewAuthTokenBlackListStorage(db *redis.Client) *AuthTokenBlackListStorage {
	return &AuthTokenBlackListStorage{
		db: db,
	}
}

func (s *AuthTokenBlackListStorage) Add(ctx context.Context, token string, ttl time.Duration) error {
	err := s.db.Set(ctx, s.key(token), "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("execute command: %w", err)
	}

	return nil
}

func (s *AuthTokenBlackListStorage) Has(ctx context.Context, token string) (bool, error) {
	if err := s.db.Get(ctx, s.key(token)).Err(); err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, fmt.Errorf("execute command: %w", err)
	}

	return true, nil
}

func (s *AuthTokenBlackListStorage) key(token string) string {
	return fmt.Sprintf("authtokenblacklist_%x", sha256.Sum256([]byte(token)))
}
