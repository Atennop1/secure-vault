package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNotFound = errors.New("no such key")
)

type Repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) Store(ctx context.Context, key string, value []byte) error {
	err := r.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("storage: failed to store key '%s' with value '%s' to redis: %w", key, value, err)
	}

	return nil
}

func (r *Repository) Load(ctx context.Context, key string) ([]byte, error) {
	value, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("storage: failed to load key '%s': %w", key, ErrNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("storage: failed to load key '%s': %w", key, err)
	}

	return value, err
}
