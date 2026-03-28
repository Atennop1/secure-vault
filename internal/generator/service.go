package generator

import (
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

func (s *Service) Generate(length int) string {
	var sb strings.Builder
	sb.Grow(length)

	for {
		for range length {
			sb.WriteByte(charset[rand.IntN(len(charset))])
		}

		if !s.repo.Contains(sb.String()) {
			break
		}

		sb.Reset()
	}

	s.repo.Store(sb.String())
	return sb.String()
}
