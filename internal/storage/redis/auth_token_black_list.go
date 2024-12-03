package redis

import (
	"context"
	"time"
)

type AuthTokenBlackListStorage struct{}

func NewAuthTokenBlackListStorage() *AuthTokenBlackListStorage {
	return &AuthTokenBlackListStorage{}
}

func (s *AuthTokenBlackListStorage) Add(ctx context.Context, token string, ttl time.Duration) error {
	return nil
}
