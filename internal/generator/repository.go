package generator

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) Store(ctx context.Context, value string) error {
	err := r.client.LPush(ctx, "slugs", value).Err()
	if err != nil {
		return fmt.Errorf("storage: failed to store slug '%s' to redis: %w", value, err)
	}

	return nil
}

func (r *Repository) Contains(ctx context.Context, value string) (bool, error) {
	_, err := r.client.LPos(ctx, "slugs", value, redis.LPosArgs{}).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, fmt.Errorf("storage: failed to check whether slug '%s' exists or not", value)
	}

	return true, nil
}
