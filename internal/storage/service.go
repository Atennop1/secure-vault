package storage

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Store(ctx context.Context, key string, value []byte) error {
	return s.repo.Store(ctx, key, value)
}

func (s *Service) Load(ctx context.Context, key string) ([]byte, error) {
	return s.repo.Load(ctx, key)
}
