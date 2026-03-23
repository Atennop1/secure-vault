package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/Atennop1/secure-vault/internal/generate"
	"github.com/Atennop1/secure-vault/internal/repository"
)

type Service struct {
	key []byte

	generator *generate.Generator
	repo      *repository.Repository
}

func NewService(key []byte, repo *repository.Repository, generator *generate.Generator) *Service {
	return &Service{
		key:       key,
		generator: generator,
		repo:      repo,
	}
}

func (s *Service) Encrypt(content string) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", fmt.Errorf("encrypt: failed to generate AES block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("encrypt: failed to generate GSM block: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("encrypt: failed to generate nonce: %w", err)
	}

	slug := s.generator.Generate(10)
	err = s.repo.Store(slug, string(gcm.Seal(nonce, nonce, []byte(content), nil)))
	if err != nil {
		return "", fmt.Errorf("encrypt: failed to store encrypted info: %w", err)
	}

	return slug, nil
}
