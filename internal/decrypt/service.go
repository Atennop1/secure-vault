package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/Atennop1/secure-vault/internal/repository"
)

type Service struct {
	key []byte

	repo *repository.Repository
}

func NewService(key []byte, repo *repository.Repository) *Service {
	return &Service{
		key:  key,
		repo: repo,
	}
}

func (s *Service) Decrypt(slug string) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", fmt.Errorf("decrypt: failed to generate AES block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("decrypt: failed to generate GSM block: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("decrypt: failed to generate nonce: %w", err)
	}

	ciphered, ok := s.repo.Load(slug)
	if !ok {
		return "", fmt.Errorf("decrypt: failed to load ciphered text from repository: %w", err)
	}

	nonceSize := gcm.NonceSize()
	data, err := gcm.Open(nil, []byte(ciphered[:nonceSize]), []byte(ciphered[nonceSize:]), nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: failed to decrypt ciphered content: %w", err)
	}

	return string(data), nil
}
