package generator

import (
	"context"
	"math/rand/v2"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Generate(ctx context.Context, length int) (string, error) {
	var sb strings.Builder
	sb.Grow(length)

	for {
		for range length {
			sb.WriteByte(charset[rand.IntN(len(charset))])
		}

		contains, err := s.repo.Contains(ctx, sb.String())
		if err != nil {
			return "", err
		}

		if !contains {
			break
		}

		sb.Reset()
	}

	err := s.repo.Store(ctx, sb.String())
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}
